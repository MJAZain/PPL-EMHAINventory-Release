package service

import (
	"context"
	"errors"
	"fmt"
	"go-gin-auth/config"
	"go-gin-auth/dto"
	"go-gin-auth/internal/adjustment"
	"go-gin-auth/internal/opname"
	"go-gin-auth/internal/product"
	"go-gin-auth/repository"
	"go-gin-auth/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockOpnameService interface {
	Create(opname *opname.StockOpname) error
	GetAll() ([]opname.StockOpname, error)
	GetByID(id uint) (opname.StockOpname, error)
	Delete(id uint) error
	IsExist(id uint) (bool, error)
	GetStockOpnameHistory(ctx context.Context) ([]dto.StockAdjustmentHistory, error)
	GetStockDiscrepancies(ctx context.Context) ([]dto.StockDiscrepancy, error)
	AdjustProductStock(ctx context.Context, productID string, req dto.StockAdjustmentRequest) (*adjustment.StockAdjustment, error)
	//Draft operations
	CreateDraft(createdBy string, opnameDate dto.DateOnly, notes string) (*opname.StockOpname, error)
	GetDraft(opnameID string) (*opname.StockOpname, error)
	UpdateDraft(opnameID string, opnameDate dto.DateOnly, notes string) (*opname.StockOpname, error)
	DeleteDraft(opnameID string) error

	AddProductToDraft(opnameID string, productID string) (*opname.StockOpnameDetail, error)
	RemoveProductFromDraft(opnameID string, detailID int) error
	// Process operations
	StartOpname(opnameID string, startedBy string) (*opname.StockOpname, error)
	RecordActualStock(detailID int, actualStock int, performedBy string, note string) (*opname.StockOpnameDetail, error)
	// Completion operations
	CompleteOpname(opnameID string, completedBy string) (*opname.StockOpname, error)
	CancelOpname(opnameID string, canceledBy string) (*opname.StockOpname, error)

	// Reporting
	GetOpnameDetails(opnameID string) (*opname.StockOpname, error)
	GetOpnameList(status opname.StockOpnameStatus, startDate, endDate time.Time, opnameID string) ([]dto.StockOpnameResponse, error)
	GetProducts(ctx context.Context) ([]dto.ProductStockResponse, error)
}

type stockOpnameService struct {
	repo repository.StockOpnameRepository
}

func NewStockOpnameService(r repository.StockOpnameRepository) StockOpnameService {
	return &stockOpnameService{r}
}

func (s *stockOpnameService) Create(opname *opname.StockOpname) error {
	db := config.DB
	for i := range opname.Details {
		d := &opname.Details[i]
		var obat product.Product
		if err := db.First(&obat, d.ProductID).Error; err != nil {
			return err
		}

		d.SystemStock = obat.StockBuffer
		d.Discrepancy = d.ActualStock - d.SystemStock

		if d.Discrepancy != 0 {
			obat.StockBuffer = d.ActualStock
			if err := db.Save(&obat).Error; err != nil {
				return err
			}
		}
	}
	return s.repo.Create(opname)
}

func (s *stockOpnameService) GetAll() ([]opname.StockOpname, error) {
	return s.repo.GetAll()
}

func (s *stockOpnameService) GetByID(id uint) (opname.StockOpname, error) {
	return s.repo.GetByID(fmt.Sprintf("%d", id))
}

func (s *stockOpnameService) Delete(id uint) error {
	return s.repo.Delete(fmt.Sprintf("%d", id))
}

func (s *stockOpnameService) IsExist(id uint) (bool, error) {
	return s.repo.IsExist(id)
}

// Fungsi untuk mendapatkan semua stok opname
func (s *stockOpnameService) GetStockOpnameHistory(ctx context.Context) ([]dto.StockAdjustmentHistory, error) {
	return s.repo.GetStockOpnameHistory(ctx)
}

func (s *stockOpnameService) GetStockDiscrepancies(ctx context.Context) ([]dto.StockDiscrepancy, error) {
	details, err := s.repo.FindAllDiscrepancies()
	if err != nil {
		return nil, err
	}

	flags, err := s.repo.FindAllFlags()
	if err != nil {
		return nil, err
	}

	var result []dto.StockDiscrepancy

	for _, detail := range details {
		detail.CalculateDiscrepancy()

		var selectedFlag string
		for _, flag := range flags {
			if detail.DiscrepancyPercentage >= flag.MinPercentage && detail.DiscrepancyPercentage <= flag.MaxPercentage {
				selectedFlag = flag.FlagName
				break
			}
		}

		result = append(result, dto.StockDiscrepancy{
			ProductID:             strconv.FormatUint(uint64(detail.ProductID), 10),
			Name:                  detail.Product.Name,
			Category:              strconv.FormatUint(uint64(detail.Product.CategoryID), 10),
			PreviousStock:         detail.SystemStock,
			ActualStock:           detail.ActualStock,
			Discrepancy:           detail.Discrepancy,
			DiscrepancyPercentage: detail.DiscrepancyPercentage,
			Flag:                  selectedFlag,
			OpnameDate:            detail.PerformedAt,
			PerformedBy:           detail.PerformedBy,
		})
	}

	return result, nil
}

func (s *stockOpnameService) AdjustProductStock(ctx context.Context, productID string, req dto.StockAdjustmentRequest) (*adjustment.StockAdjustment, error) {
	// Periksa apakah produk ada
	db := config.DB
	var product product.Product

	if err := db.First(&product, productID).Error; err != nil {
		return nil, err
	}

	// if !product.IsActive {
	// 	return nil, errors.New("cannot adjust stock for inactive product")
	// }

	// Jalankan dalam transaksi database
	var adjustment *adjustment.StockAdjustment

	err := db.Transaction(func(tx *gorm.DB) error {
		// Buat context baru dengan tx
		txCtx := context.WithValue(ctx, "tx", tx)

		// Buat adjustment
		adjustmentObj, err := s.repo.AdjustProductStock(txCtx, productID, req)
		if err != nil {
			return err
		}

		// // Simpan adjustment
		// if err := s.repo.CreateStockAdjustment(txCtx, adjustmentObj); err != nil {
		// 	return err
		// }

		// // Update stok produk
		// if err := s.repo.UpdateProductStock(txCtx, productID, req.ActualStock); err != nil {
		// 	return err
		// }

		// Update tanggal stok opname terakhir
		if err := s.repo.UpdateLastOpnameDate(txCtx, productID, req.OpnameDate); err != nil {
			return err
		}

		adjustment = adjustmentObj
		return nil
	})

	if err != nil {
		return nil, err
	}

	return adjustment, nil
}

func (s *stockOpnameService) CreateDraft(createdBy string, opnameDate dto.DateOnly, notes string) (*opname.StockOpname, error) {
	opnameID := fmt.Sprintf("OPN-%s", uuid.New().String()[:8])

	opname := &opname.StockOpname{
		OpnameID:   opnameID,
		OpnameDate: opnameDate.Local(),
		Status:     opname.Draft,
		Notes:      notes,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Jenis:      "Regular",
	}

	if err := s.repo.Create(opname); err != nil {
		return nil, err
	}

	return opname, nil
}

// GetDraft retrieves a draft stock opname by ID
func (s *stockOpnameService) GetDraft(opnameID string) (*opname.StockOpname, error) {
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return nil, err
	}

	// Only draft status can be retrieved through this method
	if data.Status != opname.Draft {
		return nil, errors.New("stock opname is not in draft status")
	}

	return &data, nil
}

// UpdateDraft updates a draft stock opname
func (s *stockOpnameService) UpdateDraft(opnameID string, opnameDate dto.DateOnly, notes string) (*opname.StockOpname, error) {
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return nil, err
	}

	if data.Status != opname.Draft {
		return nil, errors.New("only draft stock opname can be updated")
	}

	data.OpnameDate = opnameDate.Local()
	data.Notes = notes
	data.UpdatedAt = time.Now()

	if err := s.repo.Update(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteDraft deletes a draft stock opname
func (s *stockOpnameService) DeleteDraft(opnameID string) error {
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return err
	}

	if data.Status != opname.Draft {
		return errors.New("only draft stock opname can be deleted")
	}

	return s.repo.Delete(opnameID)
}

func (s *stockOpnameService) AddProductToDraft(opnameID string, productID string) (*opname.StockOpnameDetail, error) {
	// Check if opname exists and is in draft status
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return nil, err
	}

	// disable
	// if data.Status != opname.Draft {
	// 	return nil, errors.New("can only add products to draft stock opname")
	// }

	db := config.DB
	var product product.Product

	if err := db.First(&product, productID).Error; err != nil {
		return nil, err
	}
	// // Check if product exists
	// product, err := s.productRepo.FindByID(productID)
	// if err != nil {
	// 	return nil, err
	// }

	// Check if product is already in the opname
	exists, err := s.repo.ExistsByOpnameAndProduct(opnameID, productID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("product already exists in this stock opname")
	}

	productIDUint, err := utils.ConvertProductID(productID)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	// Create new detail
	detail := &opname.StockOpnameDetail{
		OpnameID:    opnameID,
		ProductID:   productIDUint,
		SystemStock: product.StockBuffer, // This needs to be retrieved from the product repository
		ActualStock: 0,                   // Will be filled during the opname process
		PerformedBy: data.CreatedBy,      // Initially set to the creator of the opname
		PerformedAt: time.Now(),
	}

	// // Calculate discrepancy for display purposes
	detail.CalculateDiscrepancy()

	if err := s.repo.CreateStockOpnameDetail(detail); err != nil {
		return nil, err
	}

	// // Setelah simpan, ambil ulang dengan relasi Product
	// if err := s.repo.PreloadProduct(detail); err != nil {
	// 	return nil, err
	// }

	// // langsung jalankan program
	// _, err = s.StartOpname(opnameID, data.CreatedBy)

	// if err != nil {
	// 	return nil, err
	// }

	return detail, nil
}

// RemoveProductFromDraft removes a product from a draft stock opname
func (s *stockOpnameService) RemoveProductFromDraft(opnameID string, detailID int) error {
	// Check if opname exists and is in draft status
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return err
	}

	if data.Status != opname.Draft {
		return errors.New("can only remove products from draft stock opname")
	}

	// Check if detail exists for this opname
	detail, err := s.repo.FindStockOpNameDetailByID(detailID)
	if err != nil {
		return err
	}

	if detail.OpnameID != opnameID {
		return errors.New("detail does not belong to this stock opname")
	}

	// Delete the detail
	return s.repo.DeleteStockOpNameDetail(detailID)
}

// StartOpname starts the stock opname process
func (s *stockOpnameService) StartOpname(opnameID string, startedBy string) (*opname.StockOpname, error) {
	// Begin transaction
	tx := config.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get opname with details
	data, err := s.repo.FindByIDWithDetails(opnameID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if data.Status != opname.Draft {
		tx.Rollback()
		return nil, errors.New("only draft stock opname can be started")
	}

	if len(data.Details) == 0 {
		tx.Rollback()
		return nil, errors.New("cannot start stock opname with no products")
	}

	// Update status to in progress
	data.Status = opname.InProgress
	data.StartTime = time.Now()
	data.UpdatedAt = time.Now()

	if err := s.repo.UpdateTx(tx, data); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return data, nil
}

// RecordActualStock records the actual stock count for a product
func (s *stockOpnameService) RecordActualStock(detailID int, actualStock int, performedBy string, note string) (*opname.StockOpnameDetail, error) {
	// Get detail
	detail, err := s.repo.FindStockOpNameDetailByID(detailID)
	if err != nil {
		return nil, err
	}

	// Get opname to check status
	// data, err := s.repo.GetByID(detail.OpnameID)
	// if err != nil {
	// 	return nil, err
	// }

	// if data.Status != opname.InProgress {
	// 	return nil, errors.New("can only record actual stock for in-progress stock opname")
	// }

	// Update detail
	detail.ActualStock = actualStock
	detail.PerformedBy = performedBy
	detail.PerformedAt = time.Now()
	detail.AdjustmentNote = note
	// Calculate discrepancy
	detail.CalculateDiscrepancy()

	if err := s.repo.UpdateStockOpNameDetail(detail); err != nil {
		return nil, err
	}

	return detail, nil
}

// CompleteOpname completes the stock opname process
func (s *stockOpnameService) CompleteOpname(opnameID string, completedBy string) (*opname.StockOpname, error) {
	// Begin transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get opname with details
	data, err := s.repo.FindByIDWithDetails(opnameID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if data.Status != opname.InProgress {
		tx.Rollback()
		return nil, errors.New("only in-progress stock opname can be completed")
	}

	// Check if all products have been counted
	for _, detail := range data.Details {
		// if detail was added but never counted (actualStock is 0)
		if detail.ActualStock == 0 && detail.AdjustmentNote == "" {
			tx.Rollback()
			return nil, errors.New("all products must be counted before completing stock opname")
		}
	}

	// Create stock adjustments for each product
	for _, detail := range data.Details {
		// Calculate discrepancy
		detail.CalculateDiscrepancy()

		// Only create adjustment if there's a discrepancy
		if detail.Discrepancy != 0 {
			adjustment := &adjustment.StockAdjustment{
				AdjustmentID:   fmt.Sprintf("ADJ-%s", uuid.New().String()[:8]),
				ProductID:      strconv.FormatUint(uint64(detail.ProductID), 10),
				PreviousStock:  detail.SystemStock,
				AdjustedStock:  detail.ActualStock,
				AdjustmentType: adjustment.Opname,
				ReferenceID:    opnameID,
				AdjustmentNote: detail.AdjustmentNote,
				AdjustmentDate: time.Now(),
				PerformedBy:    completedBy,
			}

			adjustment.CalculateAdjustmentQuantity()

			if err := s.repo.CreateStockAdjustment(tx, adjustment); err != nil {
				tx.Rollback()
				return nil, err
			}

			/// disable dulu Update product stock
			if err := s.repo.UpdateProductStock(tx, strconv.FormatUint(uint64(detail.ProductID), 10), detail.ActualStock); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Update opname status
	data.Status = opname.Completed
	data.FlagActive = false
	data.EndTime = time.Now()
	data.UpdatedAt = time.Now()

	if err := s.repo.UpdateTx(tx, data); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return data, nil
}

// CancelOpname cancels the stock opname process
func (s *stockOpnameService) CancelOpname(opnameID string, canceledBy string) (*opname.StockOpname, error) {
	// Get opname
	data, err := s.repo.GetByID(opnameID)
	if err != nil {
		return nil, err
	}

	// Only draft or in-progress can be canceled
	if data.Status != opname.Draft && data.Status != opname.InProgress {
		return nil, errors.New("only draft or in-progress stock opname can be canceled")
	}

	// Update status to canceled
	data.Status = opname.Canceled
	data.FlagActive = false
	data.EndTime = time.Now()
	data.UpdatedAt = time.Now()

	if err := s.repo.Update(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetOpnameDetails gets detailed information for a stock opname
func (s *stockOpnameService) GetOpnameDetails(opnameID string) (*opname.StockOpname, error) {
	opname, err := s.repo.FindByIDWithDetails(opnameID)
	if err != nil {
		return nil, err
	}

	// Calculate discrepancies for all details
	for i := range opname.Details {
		opname.Details[i].CalculateDiscrepancy()
	}

	return opname, nil
}

// GetOpnameList gets a list of stock opnames by status and date range
func (s *stockOpnameService) GetOpnameList(status opname.StockOpnameStatus, startDate, endDate time.Time, opnameID string) ([]dto.StockOpnameResponse, error) {
	data, err := s.repo.FindByStatusAndDateRange(status, startDate, endDate, opnameID)
	if err != nil {
		return nil, err

	}
	var responses []dto.StockOpnameResponse
	for _, o := range data {
		var detailResponses []dto.StockOpnameDetailResponse
		for _, d := range o.Details {
			detailResponses = append(detailResponses, dto.StockOpnameDetailResponse{
				ID:                     uint(d.DetailID),
				QtySystem:              d.SystemStock,
				QtyReal:                d.ActualStock,
				Discrepancy:            d.Discrepancy,
				Discrepancy_percentage: int(d.DiscrepancyPercentage),
				Adjustment_note:        d.AdjustmentNote,
				Performed_by:           d.PerformedBy,
				Performed_at:           d.PerformedAt,
				Product: dto.ProductSimple{
					ID:   d.Product.ID,
					Name: d.Product.Name,
					// Code:       d.Product.Code,
					// Barcode:    d.Product.Barcode,
					// CategoryID: d.Product.CategoryID,
				},
			})
		}

		responses = append(responses, dto.StockOpnameResponse{
			OpnameId:        o.OpnameID,
			OpnameDate:      o.OpnameDate,
			StartTime:       o.StartTime,
			EndTime:         o.EndTime,
			Status:          string(o.Status),
			Notes:           o.Notes,
			JenisStokOpname: string(o.Jenis),
			FlagActive:      o.FlagActive,
			CreatedBy:       o.CreatedBy,
			Details:         detailResponses,
		})
	}
	return responses, nil
}
func (s *stockOpnameService) GetProducts(ctx context.Context) ([]dto.ProductStockResponse, error) {
	products, err := s.repo.GetProducts(ctx)
	fmt.Println("products", products)
	if err != nil {
		return nil, err
	}
	var result []dto.ProductStockResponse
	for _, p := range products {
		result = append(result, dto.ProductStockResponse{
			Name:            p.Name,
			Code:            p.Code,
			StockBuffer:     p.StockBuffer,
			StorageLocation: p.StorageLocation.Name,
			Category: dto.CategorySimpleDTO{
				ID:   p.Category.ID,
				Name: p.Category.Name,
			},
			Unit: dto.UnitSimpleDTO{
				ID:   p.Unit.ID,
				Name: p.Unit.Name,
			},
		})
	}
	return result, nil
}
