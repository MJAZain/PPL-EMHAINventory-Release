// repository/stock_opname_repository.go
package repository

import (
	"context"
	"fmt"
	"go-gin-auth/dto"
	"go-gin-auth/internal/adjustment"
	"go-gin-auth/internal/opname"
	"go-gin-auth/internal/product"
	"go-gin-auth/model"
	"time"

	"gorm.io/gorm"
)

type StockOpnameRepository interface {
	Create(opname *opname.StockOpname) error
	GetAll() ([]opname.StockOpname, error)
	GetByID(id string) (opname.StockOpname, error)
	Delete(id string) error
	IsExist(id uint) (bool, error)
	GetStockOpnameHistory(ctx context.Context) ([]dto.StockAdjustmentHistory, error)
	GetStockDiscrepancies(ctx context.Context) ([]dto.StockDiscrepancy, error)
	AdjustProductStock(ctx context.Context, productID string, req dto.StockAdjustmentRequest) (*adjustment.StockAdjustment, error)
	CreateStockAdjustment(tx *gorm.DB, adjustment *adjustment.StockAdjustment) error
	UpdateLastOpnameDate(ctx context.Context, productID string, opnameDate time.Time) error
	UpdateProductStock(tx *gorm.DB, productID string, newStock int) error
	Update(opname *opname.StockOpname) error
	ExistsByOpnameAndProduct(opnameID string, productID string) (bool, error)
	FindByIDWithDetails(opnameID string) (*opname.StockOpname, error)
	UpdateTx(tx *gorm.DB, opname *opname.StockOpname) error

	// stock opname detail

	CreateStockOpnameDetail(detail *opname.StockOpnameDetail) error
	PreloadProduct(detail *opname.StockOpnameDetail) error
	FindStockOpNameDetailByID(detailID int) (*opname.StockOpnameDetail, error)
	UpdateStockOpNameDetail(detail *opname.StockOpnameDetail) error
	DeleteStockOpNameDetail(detailID int) error
	// reporting
	FindByStatusAndDateRange(status opname.StockOpnameStatus, startDate, endDate time.Time, opnameID string) ([]opname.StockOpname, error)
	GetProducts(ctx context.Context) ([]product.Product, error)
	FindAllFlags() ([]model.StockDiscrepancyFlag, error)
	FindAllDiscrepancies() ([]opname.StockOpnameDetail, error)
}

type stockOpnameRepository struct {
	db *gorm.DB
}

func NewStockOpnameRepository(db *gorm.DB) StockOpnameRepository {
	return &stockOpnameRepository{db}
}

func (r *stockOpnameRepository) Create(opname *opname.StockOpname) error {
	return r.db.Create(opname).Error
}

func (r *stockOpnameRepository) GetAll() ([]opname.StockOpname, error) {
	var opnames []opname.StockOpname
	err := r.db.Preload("Details").Find(&opnames).Error
	return opnames, err
}

func (r *stockOpnameRepository) GetByID(id string) (opname.StockOpname, error) {
	var opname opname.StockOpname
	err := r.db.Preload("Details").Where("opname_id = ?", id).First(&opname).Error
	return opname, err
}

func (r *stockOpnameRepository) FindAllFlags() ([]model.StockDiscrepancyFlag, error) {
	var flags []model.StockDiscrepancyFlag
	err := r.db.Find(&flags).Error
	return flags, err
}
func (r *stockOpnameRepository) Delete(id string) error {
	return r.db.Where("opname_id = ?", id).Delete(&opname.StockOpname{}).Error
}
func (r *stockOpnameRepository) IsExist(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&opname.StockOpname{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//fitur opaname

func (r *stockOpnameRepository) GetStockOpnameHistory(ctx context.Context) ([]dto.StockAdjustmentHistory, error) {
	var adjustments []adjustment.StockAdjustment

	// Ambil data stock_adjustments dengan tipe opname, urut berdasarkan tanggal desc
	err := r.db.WithContext(ctx).
		Where("adjustment_type = ?", adjustment.Opname).
		Order("adjustment_date DESC").
		Find(&adjustments).Error
	if err != nil {
		return nil, err
	}

	// Ambil nama produk untuk setiap ProductID agar bisa ditampilkan
	// Bisa ambil sekaligus nama produk semua ProductID yang unik
	productIDs := make([]string, 0, len(adjustments))
	productIDSet := make(map[string]struct{})
	for _, adj := range adjustments {
		if _, exists := productIDSet[adj.ProductID]; !exists {
			productIDs = append(productIDs, adj.ProductID)
			productIDSet[adj.ProductID] = struct{}{}
		}
	}

	var products []product.Product
	if len(productIDs) > 0 {
		err = r.db.WithContext(ctx).
			Where("id IN ?", productIDs).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
	}

	// Map ProductID ke nama produk
	productNameMap := make(map[string]string, len(products))
	for _, p := range products {
		productNameMap[fmt.Sprintf("%d", p.ID)] = p.Name
	}

	// Mapping ke DTO dan hitung discrepancy + persen di Go
	var history []dto.StockAdjustmentHistory
	for _, adj := range adjustments {
		discrepancy := adj.AdjustedStock - adj.PreviousStock

		var discrepancyPercent float64
		if adj.PreviousStock == 0 {
			if adj.AdjustedStock == 0 {
				discrepancyPercent = 0
			} else {
				discrepancyPercent = 100
			}
		} else {
			discrepancyPercent = float64(discrepancy) * 100.0 / float64(adj.PreviousStock)
		}

		history = append(history, dto.StockAdjustmentHistory{
			AdjustmentID:          adj.AdjustmentID,
			ProductID:             adj.ProductID,
			Name:                  productNameMap[adj.ProductID],
			PreviousStock:         adj.PreviousStock,
			ActualStock:           adj.AdjustedStock,
			Discrepancy:           discrepancy,
			DiscrepancyPercentage: discrepancyPercent,
			AdjustmentNote:        adj.AdjustmentNote,
			OpnameDate:            adj.AdjustmentDate,
			PerformedBy:           adj.PerformedBy,
		})
	}

	return history, nil
}

func (r *stockOpnameRepository) GetStockOpnameHistoryOld(ctx context.Context) ([]dto.StockAdjustmentHistory, error) {
	var history []dto.StockAdjustmentHistory

	result := r.db.WithContext(ctx).
		Table("stock_adjustments").
		Select(
			"stock_adjustments.adjustment_id, stock_adjustments.product_id, "+
				"products.name, stock_adjustments.previous_stock, "+
				"stock_adjustments.adjusted_stock as actual_stock, "+
				"(stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) as discrepancy, "+
				"CASE WHEN stock_adjustments.previous_stock = 0 THEN "+
				"  CASE WHEN stock_adjustments.adjusted_stock = 0 THEN 0 ELSE 100 END "+
				"ELSE "+
				"  ((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / stock_adjustments.previous_stock) "+
				"END as discrepancy_percentage, "+
				"stock_adjustments.adjustment_note, stock_adjustments.adjustment_date as opname_date, "+
				"stock_adjustments.performed_by",
		).
		Joins("JOIN products ON stock_adjustments.product_id = products.id").
		Where("stock_adjustments.adjustment_type = ?", adjustment.Opname).
		Order("stock_adjustments.adjustment_date DESC").
		Scan(&history)

	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func (r *stockOpnameRepository) GetStockDiscrepancies(ctx context.Context) ([]dto.StockDiscrepancy, error) {
	var discrepancies []dto.StockDiscrepancy

	// Query untuk mendapatkan selisih stok yang signifikan (lebih dari 10%)
	result := r.db.WithContext(ctx).
		Table("stock_adjustments").
		Select(
			"stock_adjustments.product_id, products.name, products.category, "+
				"stock_adjustments.previous_stock, stock_adjustments.adjusted_stock as actual_stock, "+
				"(stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) as discrepancy, "+
				"CASE WHEN stock_adjustments.previous_stock = 0 THEN "+
				"  CASE WHEN stock_adjustments.adjusted_stock = 0 THEN 0 ELSE 100 END "+
				"ELSE "+
				"  ((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / stock_adjustments.previous_stock) "+
				"END as discrepancy_percentage, "+
				"CASE "+
				"  WHEN ABS(((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / NULLIF(stock_adjustments.previous_stock, 0))) > 20 THEN 'HIGH_LOSS' "+
				"  WHEN ((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / NULLIF(stock_adjustments.previous_stock, 0)) < -10 THEN 'MODERATE_LOSS' "+
				"  WHEN ((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / NULLIF(stock_adjustments.previous_stock, 0)) > 10 THEN 'HIGH_GAIN' "+
				"  ELSE 'NORMAL' "+
				"END as flag, "+
				"stock_adjustments.adjustment_date as opname_date, "+
				"stock_adjustments.performed_by",
		).
		Joins("JOIN products ON stock_adjustments.product_id = products.product_id").
		Where("stock_adjustments.adjustment_type = ?", adjustment.Opname).
		Where(
			"ABS(((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / NULLIF(stock_adjustments.previous_stock, 0))) > 10",
		).
		Order("ABS(((stock_adjustments.adjusted_stock - stock_adjustments.previous_stock) * 100.0 / NULLIF(stock_adjustments.previous_stock, 0))) DESC").
		Scan(&discrepancies)

	if result.Error != nil {
		return nil, result.Error
	}

	return discrepancies, nil
}

// func (r *stockOpnameRepository) FindAllDiscrepancies() ([]opname.StockOpnameDetail, error) {
// 	var details []opname.StockOpnameDetail

// 	err := r.db.
// 		Joins("JOIN stock_opnames ON stock_opnames.opname_id = stock_opname_details.opname_id").
// 		Where("stock_opnames.status = ?", "completed").
// 		Find(&details).Error

// 	return details, err
// }

func (r *stockOpnameRepository) FindAllDiscrepancies() ([]opname.StockOpnameDetail, error) {
	var details []opname.StockOpnameDetail

	err := r.db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "code", "barcode", "category_id")
		}).
		Joins("JOIN stock_opnames ON stock_opnames.opname_id = stock_opname_details.opname_id").
		Where("stock_opnames.status = ?", "completed").
		Find(&details).Error

	return details, err
}

func (r *stockOpnameRepository) AdjustProductStock(ctx context.Context, productID string, req dto.StockAdjustmentRequest) (*adjustment.StockAdjustment, error) {
	// Dapatkan stok saat ini
	var productStock model.ProductStock
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&productStock).Error; err != nil {
		return nil, err
	}

	// Buat penyesuaian stok
	adjustment := adjustment.StockAdjustment{
		AdjustmentID:   generateAdjustmentID(), // Implementasi fungsi ini sesuai kebutuhan
		ProductID:      productID,
		PreviousStock:  productStock.CurrentStock,
		AdjustedStock:  req.ActualStock,
		AdjustmentType: adjustment.Opname,
		ReferenceID:    "", // Bisa diisi dengan ID opname jika perlu
		AdjustmentNote: req.AdjustmentNote,
		AdjustmentDate: req.OpnameDate,
		PerformedBy:    req.PerformedBy,
	}

	adjustment.CalculateAdjustmentQuantity()

	return &adjustment, nil
}

//	func (r *stockOpnameRepository) CreateStockAdjustment(ctx context.Context, adjustment model.StockAdjustment) error {
//		return r.db.WithContext(ctx).Create(&adjustment).Error
//	}
func (r *stockOpnameRepository) CreateStockAdjustment(tx *gorm.DB, adjustment *adjustment.StockAdjustment) error {
	return tx.Create(adjustment).Error
}

func (r *stockOpnameRepository) UpdateLastOpnameDate(ctx context.Context, productID string, opnameDate time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&model.ProductStock{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"last_opname_date": opnameDate,
		})

	return result.Error
}

// Helper function untuk generate ID penyesuaian
func generateAdjustmentID() string {
	// Implementasi sesuai kebutuhan, misalnya:
	return "ADJ" + time.Now().Format("20060102150405")
}

// func (r *stockOpnameRepository) UpdateProductStock(tx *gorm.DB, productID string, newStock int) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.ProductStock{}).
// 		Where("product_id = ?", productID).
// 		Updates(map[string]interface{}{
// 			"current_stock": newStock,
// 		})

//		return result.Error
//	}
func (r *stockOpnameRepository) UpdateProductStock(tx *gorm.DB, productID string, newStock int) error {
	result := tx.
		Model(&product.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"stock_buffer": newStock,
		})

	return result.Error
}

func (r *stockOpnameRepository) Update(opname *opname.StockOpname) error {
	return r.db.Save(opname).Error
}
func (r *stockOpnameRepository) ExistsByOpnameAndProduct(opnameID string, productID string) (bool, error) {
	var count int64
	if err := r.db.Model(&opname.StockOpnameDetail{}).
		Where("opname_id = ? AND product_id = ?", opnameID, productID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *stockOpnameRepository) CreateStockOpnameDetail(detail *opname.StockOpnameDetail) error {
	return r.db.Create(detail).Error
}
func (r *stockOpnameRepository) PreloadProduct(detail *opname.StockOpnameDetail) error {
	return r.db.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(detail, detail.DetailID).Error
}

func (r *stockOpnameRepository) FindByIDWithDetails(opnameID string) (*opname.StockOpname, error) {
	var opname opname.StockOpname
	if err := r.db.Preload("Details").Where("opname_id = ?", opnameID).First(&opname).Error; err != nil {
		return nil, err
	}
	return &opname, nil
}
func (r *stockOpnameRepository) UpdateTx(tx *gorm.DB, opname *opname.StockOpname) error {
	return tx.Save(opname).Error
}
func (r *stockOpnameRepository) FindStockOpNameDetailByID(detailID int) (*opname.StockOpnameDetail, error) {
	var detail opname.StockOpnameDetail
	if err := r.db.Where("detail_id = ?", detailID).First(&detail).Error; err != nil {
		return nil, err
	}
	return &detail, nil
}
func (r *stockOpnameRepository) UpdateStockOpNameDetail(detail *opname.StockOpnameDetail) error {
	return r.db.Save(detail).Error
}
func (r *stockOpnameRepository) DeleteStockOpNameDetail(detailID int) error {
	return r.db.Where("detail_id = ?", detailID).Delete(&opname.StockOpnameDetail{}).Error
}
func (r *stockOpnameRepository) FindByStatusAndDateRange(status opname.StockOpnameStatus, startDate, endDate time.Time, opnameID string) ([]opname.StockOpname, error) {
	var opnames []opname.StockOpname
	query := r.db

	if opnameID != "" {
		query = query.Where("opname_id = ?", opnameID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("opname_date BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("opname_date >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("opname_date <= ?", endDate)
	}

	// if err := query.Order("opname_date DESC").Find(&opnames).Error; err != nil {
	// 	return nil, err
	// }

	// 👇 Ini kuncinya!
	if err := query.
		Preload("Details").
		Preload("Details.Product").
		Order("opname_date DESC").
		Find(&opnames).Error; err != nil {
		return nil, err
	}
	return opnames, nil
}
func (r *stockOpnameRepository) GetProducts(ctx context.Context) ([]product.Product, error) {
	var products []product.Product
	err := r.db.WithContext(ctx).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
