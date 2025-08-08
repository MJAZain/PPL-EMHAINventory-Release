package sales

import (
	"errors"
	"fmt"
	"go-gin-auth/internal/stock"
	"time"

	"gorm.io/gorm"
)

// Service interface
type SalesRegularService interface {
	GetAll(limit, offset int) ([]SalesRegular, int64, error)
	GetByID(id uint) (*SalesRegular, error)
	Create(req *SalesRegularRequest) (*SalesRegular, error)
	Update(id uint, req *SalesRegularRequest) (*SalesRegular, error)
	Delete(id uint) error
}

// Service struct
type salesRegularService struct {
	db        *gorm.DB
	repo      SalesRegularRepository
	stockRepo stock.Repository
}

// Constructor
func NewSalesRegularService(db *gorm.DB, repo SalesRegularRepository, stockRepo stock.Repository) *salesRegularService {
	return &salesRegularService{
		db:        db,
		repo:      repo,
		stockRepo: stockRepo,
	}
}

// ✅ Ambil semua penjualan dengan paginasi
func (s *salesRegularService) GetAll(limit, offset int) ([]SalesRegular, int64, error) {
	return s.repo.GetAllSalesRegular(limit, offset)
}

// ✅ Ambil detail transaksi
func (s *salesRegularService) GetByID(id uint) (*SalesRegular, error) {
	return s.repo.GetSalesRegularByID(id)
}

// ✅ Create new sales regular
func (s *salesRegularService) Create(req *SalesRegularRequest) (*SalesRegular, error) {
	tx := s.db.Begin()

	salesCode := fmt.Sprintf("SR-%d", time.Now().UnixNano())

	newSale := &SalesRegular{
		SalesCode:       salesCode,
		TransactionDate: req.TransactionDate,
		CashierName:     req.CashierName,
		CustomerName:    req.CustomerName,
		CustomerContact: req.CustomerContact,
		Description:     req.Description,
		SubTotal:        req.SubTotal,
		TotalDiscount:   req.TotalDiscount,
		TotalPay:        req.TotalPay,
		PaymentMethod:   req.PaymentMethod,
		ShiftID:         req.ShiftID,
	}

	if err := tx.Create(newSale).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range req.Items {
		// Kurangi stok
		if err := s.stockRepo.UpdateProductStock(item.ProductID, item.Qty, false); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal mengurangi stok: %w", err)
		}

		newItem := SalesRegularItem{
			SalesRegularID: newSale.ID,
			ProductID:      item.ProductID,
			ProductCode:    item.ProductCode,
			ProductName:    item.ProductName,
			Qty:            item.Qty,
			Unit:           item.Unit,
			UnitPrice:      item.UnitPrice,
			SubTotal:       item.SubTotal,
		}

		if err := tx.Create(&newItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newSale, nil
}

func (s *salesRegularService) Update(id uint, req *SalesRegularRequest) (*SalesRegular, error) {
	tx := s.db.Begin()

	// Ambil transaksi beserta itemnya dengan preload
	existing := &SalesRegular{}
	if err := tx.Preload("Items").First(existing, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("transaksi tidak ditemukan")
	}

	// Step 1: Kembalikan stok lama (rollback stok ke stok semula)
	for _, oldItem := range existing.Items {
		if err := s.stockRepo.UpdateProductStock(oldItem.ProductID, oldItem.Qty, true); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal mengembalikan stok lama: %w", err)
		}
	}

	// Step 2: Soft delete semua item lama
	if err := tx.Where("sales_regular_id = ?", id).Delete(&SalesRegularItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Step 3: Tambah item baru dan kurangi stok
	for _, item := range req.Items {
		if err := s.stockRepo.UpdateProductStock(item.ProductID, item.Qty, false); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("stok tidak mencukupi untuk produk %d", item.ProductID)
		}

		newItem := SalesRegularItem{
			SalesRegularID: id,
			ProductID:      item.ProductID,
			ProductCode:    item.ProductCode,
			ProductName:    item.ProductName,
			Qty:            item.Qty,
			Unit:           item.Unit,
			UnitPrice:      item.UnitPrice,
			SubTotal:       item.SubTotal,
		}
		if err := tx.Create(&newItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Step 4: Update header transaksi langsung di tx
	existing.TransactionDate = req.TransactionDate
	existing.CashierName = req.CashierName
	existing.CustomerName = req.CustomerName
	existing.CustomerContact = req.CustomerContact
	existing.Description = req.Description
	existing.SubTotal = req.SubTotal
	existing.TotalDiscount = req.TotalDiscount
	existing.TotalPay = req.TotalPay
	existing.PaymentMethod = req.PaymentMethod
	existing.ShiftID = req.ShiftID
	existing.UpdatedAt = time.Now()

	if err := tx.Save(existing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *salesRegularService) Delete(id uint) error {
	tx := s.db.Begin()

	// Ambil transaksi beserta itemnya
	sale, err := s.repo.GetSalesRegularByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Kembalikan stok semua item yang ada
	for _, item := range sale.Items {
		if err := s.stockRepo.UpdateProductStock(item.ProductID, item.Qty, true); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Soft delete item terkait
	if err := tx.Where("sales_regular_id = ?", id).Delete(&SalesRegularItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Soft delete transaksi penjualan
	if err := tx.Delete(&SalesRegular{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
