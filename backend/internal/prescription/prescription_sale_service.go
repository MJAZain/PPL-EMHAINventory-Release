package prescription

import (
	"fmt"
	"go-gin-auth/internal/stock"
	"log"

	"gorm.io/gorm"
)

// Service
type PrescriptionSaleService struct {
	db *gorm.DB
}

func NewPrescriptionSaleService(db *gorm.DB) *PrescriptionSaleService {
	return &PrescriptionSaleService{db: db}
}

func (s *PrescriptionSaleService) GetAll(page, limit int) ([]PrescriptionSale, int64, error) {
	var sales []PrescriptionSale
	var total int64

	offset := (page - 1) * limit

	err := s.db.Model(&PrescriptionSale{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Preload("Doctor").Preload("Patient").Preload("Shift").
		Preload("Items").Preload("Items.Stock").
		Offset(offset).Limit(limit).
		Order("created_at DESC").Find(&sales).Error

	return sales, total, err
}

// GetByID retrieves a prescription sale by ID with all related data
func (s *PrescriptionSaleService) GetByID(id uint) (*PrescriptionSale, error) {
	var sale PrescriptionSale

	// Use raw query to avoid any preload issues that might cause duplicates
	err := s.db.Where("id = ?", id).First(&sale).Error
	if err != nil {
		return nil, err
	}

	// Load relations separately to avoid complex joins
	s.db.Where("id = ?", sale.DoctorID).First(&sale.Doctor)
	s.db.Where("id = ?", sale.PatientID).First(&sale.Patient)
	s.db.Where("id = ?", sale.ShiftID).First(&sale.Shift)

	// Load items with explicit query to avoid duplicates
	var items []PrescriptionItem
	s.db.Where("prescription_sale_id = ?", id).Find(&items)

	// Load stock for each item
	for i := range items {
		s.db.Where("id = ?", items[i].StockID).First(&items[i].Stock)
	}

	sale.Items = items

	return &sale, nil
}

func (s *PrescriptionSaleService) Create(req *CreatePrescriptionSaleRequest) (*PrescriptionSale, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate stock availability first
	if err := s.validateStockAvailability(tx, req.Items); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Generate transaction code
	transactionCode := generateTransactionCode()

	// Calculate total
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Apply discount
	if req.DiscountPercent > 0 {
		totalAmount = totalAmount * (1 - req.DiscountPercent/100)
	}
	totalAmount -= req.DiscountAmount

	// Create prescription sale
	sale := PrescriptionSale{
		TransactionCode:  transactionCode,
		PrescriptionNo:   req.PrescriptionNo,
		PrescriptionDate: req.PrescriptionDate,
		DoctorID:         req.DoctorID,
		Clinic:           req.Clinic,
		Diagnosis:        req.Diagnosis,
		PatientID:        req.PatientID,
		TransactionDate:  req.TransactionDate,
		PaymentMethod:    req.PaymentMethod,
		DiscountPercent:  req.DiscountPercent,
		DiscountAmount:   req.DiscountAmount,
		TotalAmount:      totalAmount,
		ShiftID:          req.ShiftID,
	}

	if err := tx.Create(&sale).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create prescription sale: %w", err)
	}

	// Create items and update stock
	for _, itemReq := range req.Items {
		// Get stock berdasarkan ProductID
		var stockItem stock.Stock
		err := tx.Where("product_id = ?", itemReq.ProductID).First(&stockItem).Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("stock not found for product ID %d: %w", itemReq.ProductID, err)
		}

		// Create prescription item
		item := PrescriptionItem{
			PrescriptionSaleID: sale.ID,
			StockID:            stockItem.ID,
			ItemCode:           itemReq.Code,
			ItemName:           itemReq.Name,
			Quantity:           itemReq.Quantity,
			Unit:               itemReq.Unit,
			Price:              itemReq.Price,
			SubTotal:           itemReq.Price * float64(itemReq.Quantity),
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create prescription item: %w", err)
		}

		// Update stock - MENGURANGI stock karena barang terjual
		result := tx.Model(&stock.Stock{}).
			Where("id = ? AND quantity >= ?", stockItem.ID, itemReq.Quantity).
			Update("quantity", gorm.Expr("quantity - ?", itemReq.Quantity))

		if result.Error != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock quantity: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product ID %d", itemReq.ProductID)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return with preloaded data
	return s.GetByID(sale.ID)
}

func (s *PrescriptionSaleService) Update(id uint, req *CreatePrescriptionSaleRequest) (*PrescriptionSale, error) {
	log.Printf("🔄 [Start] Updating prescription sale ID %d", id)

	// Step 1: Get existing sale with items
	existingSale, err := s.GetByID(id)
	if err != nil {
		log.Printf("❌ Failed to get existing sale ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get existing sale: %w", err)
	}

	// Step 2: Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		log.Printf("❌ Failed to start transaction: %v", tx.Error)
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	log.Println("✅ Transaction started")

	// Step 3: Calculate net stock changes
	log.Println("📊 Calculating stock changes...")
	productStockChanges := make(map[uint]int)

	for _, itemReq := range req.Items {
		productStockChanges[itemReq.ProductID] += itemReq.Quantity
	}
	for _, existingItem := range existingSale.Items {
		var existingStock stock.Stock
		if err := tx.Where("id = ?", existingItem.StockID).First(&existingStock).Error; err != nil {
			log.Printf("❌ Failed to get stock info: %v", err)
			tx.Rollback()
			return nil, err
		}
		productStockChanges[existingStock.ProductID] -= existingItem.Quantity
	}

	// Step 4: Validate stock availability
	for productID, netChange := range productStockChanges {
		if netChange > 0 {
			var currentStock stock.Stock
			if err := tx.Where("product_id = ?", productID).First(&currentStock).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("stock not found for product ID %d: %w", productID, err)
			}
			if currentStock.Quantity < netChange {
				tx.Rollback()
				return nil, fmt.Errorf("insufficient stock for product ID %d", productID)
			}
		}
	}

	// Step 5: Restore stock from existing items
	for _, item := range existingSale.Items {
		if err := tx.Model(&stock.Stock{}).
			Where("id = ?", item.StockID).
			Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to restore stock: %w", err)
		}
	}

	// Step 6: Soft-delete existing items
	for _, item := range existingSale.Items {
		if err := tx.Delete(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete item ID %d: %w", item.ID, err)
		}
	}

	// Step 7: Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	if req.DiscountPercent > 0 {
		totalAmount *= (1 - req.DiscountPercent/100)
	}
	totalAmount -= req.DiscountAmount

	// Step 8: Update sale main record
	updates := map[string]interface{}{
		"prescription_no":   req.PrescriptionNo,
		"prescription_date": req.PrescriptionDate,
		"doctor_id":         req.DoctorID,
		"clinic":            req.Clinic,
		"diagnosis":         req.Diagnosis,
		"patient_id":        req.PatientID,
		"transaction_date":  req.TransactionDate,
		"payment_method":    req.PaymentMethod,
		"discount_percent":  req.DiscountPercent,
		"discount_amount":   req.DiscountAmount,
		"total_amount":      totalAmount,
		"shift_id":          req.ShiftID,
	}
	if err := tx.Model(&existingSale).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update prescription sale: %w", err)
	}

	// Step 9: Create new items and update stock
	for _, itemReq := range req.Items {
		var stockItem stock.Stock
		if err := tx.Where("product_id = ?", itemReq.ProductID).First(&stockItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("stock not found for product ID %d: %w", itemReq.ProductID, err)
		}

		item := PrescriptionItem{
			PrescriptionSaleID: id,
			StockID:            stockItem.ID,
			ItemCode:           itemReq.Code,
			ItemName:           itemReq.Name,
			Quantity:           itemReq.Quantity,
			Unit:               itemReq.Unit,
			Price:              itemReq.Price,
			SubTotal:           itemReq.Price * float64(itemReq.Quantity),
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create prescription item: %w", err)
		}

		result := tx.Model(&stock.Stock{}).
			Where("id = ? AND quantity >= ?", stockItem.ID, itemReq.Quantity).
			Update("quantity", gorm.Expr("quantity - ?", itemReq.Quantity))

		if result.Error != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product ID %d", itemReq.ProductID)
		}
	}

	// Step 10: Final verification
	var finalItemCount int64
	if err := tx.Model(&PrescriptionItem{}).
		Where("prescription_sale_id = ?", id).
		Count(&finalItemCount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to count final items: %w", err)
	}
	if finalItemCount != int64(len(req.Items)) {
		tx.Rollback()
		return nil, fmt.Errorf("data integrity error: expected %d items but found %d", len(req.Items), finalItemCount)
	}

	// Step 11: Commit
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit update: %w", err)
	}

	// Step 12: Return updated sale
	return s.GetByID(id)
}

func (s *PrescriptionSaleService) Delete(id uint) error {
	// Get existing sale with items
	sale, err := s.GetByID(id)
	if err != nil {
		return err
	}

	tx := s.db.Begin()

	// Restore stock - MENAMBAH stock karena barang dikembalikan
	for _, item := range sale.Items {
		if err := tx.Model(&stock.Stock{}).Where("id = ?", item.StockID).
			Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error; err != nil {
			return fmt.Errorf("failed to restore stock: %w", err)
		}
	}

	// Delete prescription items first (foreign key constraint)
	if err := tx.Where("prescription_sale_id = ?", id).Delete(&PrescriptionItem{}).Error; err != nil {
		return fmt.Errorf("failed to delete prescription items: %w", err)
	}

	// Delete prescription sale
	if err := tx.Delete(&PrescriptionSale{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete prescription sale: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper function to validate stock availability
func (s *PrescriptionSaleService) validateStockAvailability(tx *gorm.DB, items []CreatePrescriptionItemRequest) error {
	for _, item := range items {
		var stockItem stock.Stock
		err := tx.Where("product_id = ?", item.ProductID).First(&stockItem).Error
		if err != nil {
			return fmt.Errorf("stock not found for product ID %d: %w", item.ProductID, err)
		}

		if stockItem.Quantity < item.Quantity {
			return fmt.Errorf("insufficient stock for product ID %d. Available: %d, Required: %d",
				item.ProductID, stockItem.Quantity, item.Quantity)
		}
	}
	return nil
}
