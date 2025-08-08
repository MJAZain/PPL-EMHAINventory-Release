package pbf

import (
	"fmt"
	"go-gin-auth/config"
	"go-gin-auth/internal/product"
	"go-gin-auth/internal/stock"
	"go-gin-auth/internal/unit"
	"go-gin-auth/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handlers
func GetAllIncomingPBF(c *gin.Context) {
	var incomingPBFs []IncomingPBF

	// Get query parameters for pagination and filtering
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := config.DB.Preload("Details.Product").Preload("Supplier").Preload("User")

	// Filter by supplier if provided
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}

	// Filter by date range if provided
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("receipt_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("receipt_date <= ?", endDate)
	}

	result := query.Offset(offset).Limit(limit).Find(&incomingPBFs)
	if result.Error != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve data", result.Error.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data retrieved successfully", nil, incomingPBFs)
}
func CreateIncomingPBF(c *gin.Context) {
	var req CreateIncomingPBFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request body", err.Error(), nil)
		return
	}

	// Parse dates
	orderDate, err := parseDate(req.OrderDate)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid order date format", err.Error(), nil)
		return
	}

	receiptDate, err := parseDate(req.ReceiptDate)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid receipt date format", err.Error(), nil)
		return
	}

	// Parse payment due date if provided
	var paymentDueDate *time.Time
	if req.PaymentDueDate != nil && *req.PaymentDueDate != "" {
		parsedDate, err := parseDate(*req.PaymentDueDate)
		if err != nil {
			utils.Respond(c, http.StatusBadRequest, "Invalid payment due date format", err.Error(), nil)
			return
		}
		paymentDueDate = &parsedDate
	}

	// **BEGIN TRANSACTION**
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Calculate total purchase
	var totalPurchase float64
	var details []IncomingPBFDetail

	for _, detailReq := range req.Details {
		// Get product info
		var product product.Product
		// Ambil data produk
		if err := config.DB.First(&product, detailReq.ProductID).Error; err != nil {
			utils.Respond(c, http.StatusBadRequest, fmt.Sprintf("Product with ID %d not found", detailReq.ProductID), err.Error(), nil)
			return
		}

		// Ambil data unit secara manual
		var unit unit.Unit
		if err := config.DB.First(&unit, product.UnitID).Error; err == nil {
			product.Unit = unit
		}

		totalPrice := float64(detailReq.Quantity) * detailReq.PurchasePrice
		totalPurchase += totalPrice

		// Parse expiry date if provided
		var expiryDate *time.Time
		if detailReq.ExpiryDate != nil && *detailReq.ExpiryDate != "" {
			parsedDate, err := parseDate(*detailReq.ExpiryDate)
			if err != nil {
				utils.Respond(c, http.StatusBadRequest, "Invalid expiry date format", err.Error(), nil)
				return
			}
			expiryDate = &parsedDate
		}

		detail := IncomingPBFDetail{
			ProductID:     detailReq.ProductID,
			ProductCode:   product.Code,
			ProductName:   product.Name,
			Unit:          unit.Name,
			Quantity:      detailReq.Quantity,
			PurchasePrice: detailReq.PurchasePrice,
			TotalPrice:    totalPrice,
			BatchNumber:   detailReq.BatchNumber,
			ExpiryDate:    expiryDate,
		}
		details = append(details, detail)
	}

	// Set default payment status if not provided
	paymentStatus := req.PaymentStatus
	if paymentStatus == "" {
		paymentStatus = "Belum Lunas"
	}

	// Create incoming PBF record
	incomingPBF := IncomingPBF{
		OrderNumber:     req.OrderNumber,
		OrderDate:       orderDate,
		ReceiptDate:     receiptDate,
		TransactionCode: generateTransactionCode(),
		SupplierID:      req.SupplierID,
		InvoiceNumber:   req.InvoiceNumber,
		TransactionType: req.TransactionType,
		PaymentDueDate:  paymentDueDate,
		UserID:          req.UserID,
		AdditionalNotes: req.AdditionalNotes,
		TotalPurchase:   totalPurchase,
		PaymentStatus:   paymentStatus,
		Details:         details,
	}

	// Save to database
	if err := config.DB.Create(&incomingPBF).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to create incoming PBF record", err.Error(), nil)
		return
	}

	// **UPDATE STOCK FOR EACH DETAIL**
	for _, detail := range details {
		if err := updateStock(tx, detail, "ADD"); err != nil {
			tx.Rollback()
			utils.Respond(c, http.StatusInternalServerError, "Failed to update stock", err.Error(), nil)
			return
		}
	}

	// **COMMIT TRANSACTION**
	if err := tx.Commit().Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to commit transaction", err.Error(), nil)
		return
	}

	// Load relations for response
	config.DB.Preload("Details.Product").Preload("Supplier").Preload("User").First(&incomingPBF, incomingPBF.ID)

	utils.Respond(c, http.StatusCreated, "Incoming PBF record created successfully", nil, incomingPBF)
}

func GetIncomingPBFByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID", err.Error(), nil)
		return
	}

	var incomingPBF IncomingPBF
	result := config.DB.Preload("Details.Product").Preload("Supplier").Preload("User").First(&incomingPBF, id)
	if result.Error != nil {
		utils.Respond(c, http.StatusNotFound, "Data not found", result.Error.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data retrieved successfully", nil, incomingPBF)
}

func UpdateIncomingPBF(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID", err.Error(), nil)
		return
	}

	// **BEGIN TRANSACTION**
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if record exists
	var existingRecord IncomingPBF
	if err := config.DB.First(&existingRecord, id).Error; err != nil {
		utils.Respond(c, http.StatusNotFound, "Record not found", err.Error(), nil)
		return
	}
	// **GET OLD DETAILS FOR STOCK REVERSAL**
	var oldDetails []IncomingPBFDetail
	if err := tx.Where("incoming_pbf_id = ?", id).Find(&oldDetails).Error; err != nil {
		tx.Rollback()
		utils.Respond(c, http.StatusInternalServerError, "Failed to get old details", err.Error(), nil)
		return
	}

	// **REVERT OLD STOCK CHANGES**
	for _, oldDetail := range oldDetails {
		if err := updateStock(tx, oldDetail, "SUBTRACT"); err != nil {
			tx.Rollback()
			utils.Respond(c, http.StatusInternalServerError, "Failed to revert stock", err.Error(), nil)
			return
		}
	}
	var req CreateIncomingPBFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request body", err.Error(), nil)
		return
	}

	// Delete existing details
	config.DB.Where("incoming_pbf_id = ?", id).Delete(&IncomingPBFDetail{})

	// Parse dates and create updated record (similar to create function)
	orderDate, _ := parseDate(req.OrderDate)
	receiptDate, _ := parseDate(req.ReceiptDate)

	var paymentDueDate *time.Time
	if req.PaymentDueDate != nil && *req.PaymentDueDate != "" {
		parsedDate, _ := parseDate(*req.PaymentDueDate)
		paymentDueDate = &parsedDate
	}

	var totalPurchase float64
	var details []IncomingPBFDetail

	for _, detailReq := range req.Details {
		var product product.Product
		// Ambil data produk
		if err := config.DB.First(&product, detailReq.ProductID).Error; err != nil {
			utils.Respond(c, http.StatusBadRequest,
				fmt.Sprintf("Product with ID %d not found", detailReq.ProductID),
				err.Error(), nil)
			return
		}

		// Ambil data unit secara manual
		var unit unit.Unit
		if err := config.DB.First(&unit, product.UnitID).Error; err == nil {
			product.Unit = unit
		}

		totalPrice := float64(detailReq.Quantity) * detailReq.PurchasePrice
		totalPurchase += totalPrice

		var expiryDate *time.Time
		if detailReq.ExpiryDate != nil && *detailReq.ExpiryDate != "" {
			parsedDate, _ := parseDate(*detailReq.ExpiryDate)
			expiryDate = &parsedDate
		}

		detail := IncomingPBFDetail{
			IncomingPBFID: uint(id),
			ProductID:     detailReq.ProductID,
			ProductCode:   product.Code,
			ProductName:   product.Name,
			Unit:          unit.Name,
			Quantity:      detailReq.Quantity,
			PurchasePrice: detailReq.PurchasePrice,
			TotalPrice:    totalPrice,
			BatchNumber:   detailReq.BatchNumber,
			ExpiryDate:    expiryDate,
		}
		details = append(details, detail)
	}

	// Update main record
	updates := map[string]interface{}{
		"order_number":     req.OrderNumber,
		"order_date":       orderDate,
		"receipt_date":     receiptDate,
		"supplier_id":      req.SupplierID,
		"invoice_number":   req.InvoiceNumber,
		"transaction_type": req.TransactionType,
		"payment_due_date": paymentDueDate,
		"user_id":          req.UserID,
		"additional_notes": req.AdditionalNotes,
		"total_purchase":   totalPurchase,
		"payment_status":   req.PaymentStatus,
	}

	if err := config.DB.Model(&existingRecord).Updates(updates).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to update record", err.Error(), nil)
		return
	}

	// Create new details
	// for _, detail := range details {
	// 	config.DB.Create(&detail)
	// }

	// Create new details and update stock
	for _, detail := range details {
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			utils.Respond(c, http.StatusInternalServerError, "Failed to create detail", err.Error(), nil)
			return
		}

		// **UPDATE STOCK FOR NEW DETAILS**
		if err := updateStock(tx, detail, "ADD"); err != nil {
			tx.Rollback()
			utils.Respond(c, http.StatusInternalServerError, "Failed to update stock", err.Error(), nil)
			return
		}
	}

	// **COMMIT TRANSACTION**
	if err := tx.Commit().Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to commit transaction", err.Error(), nil)
		return
	}

	// Load updated record with relations
	var updatedRecord IncomingPBF
	config.DB.Preload("Details.Product").Preload("Supplier").Preload("User").First(&updatedRecord, id)

	utils.Respond(c, http.StatusOK, "Record updated successfully", nil, updatedRecord)
}

func DeleteIncomingPBF(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID", err.Error(), nil)
		return
	}
	// **BEGIN TRANSACTION**
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if record exists
	var existingRecord IncomingPBF
	if err := config.DB.First(&existingRecord, id).Error; err != nil {
		utils.Respond(c, http.StatusNotFound, "Record not found", err.Error(), nil)
		return
	}
	// **GET DETAILS FOR STOCK REVERSAL**
	var details []IncomingPBFDetail
	if err := tx.Where("incoming_pbf_id = ?", id).Find(&details).Error; err != nil {
		tx.Rollback()
		utils.Respond(c, http.StatusInternalServerError, "Failed to get details", err.Error(), nil)
		return
	}

	// **REVERT STOCK CHANGES**
	for _, detail := range details {
		if err := updateStock(tx, detail, "SUBTRACT"); err != nil {
			tx.Rollback()
			utils.Respond(c, http.StatusInternalServerError, "Failed to revert stock", err.Error(), nil)
			return
		}
	}

	// Delete details first (due to foreign key constraint)
	config.DB.Where("incoming_pbf_id = ?", id).Delete(&IncomingPBFDetail{})

	// Delete main record
	if err := config.DB.Delete(&existingRecord).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to delete record", err.Error(), nil)
		return
	}
	// **COMMIT TRANSACTION**
	if err := tx.Commit().Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to commit transaction", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Record deleted successfully", nil, nil)
}

// **FUNGSI UTAMA UNTUK UPDATE STOCK**
func updateStock(tx *gorm.DB, detail IncomingPBFDetail, operation string) error {
	var stockdata stock.Stock

	err := tx.Where("product_id = ?", detail.ProductID).
		First(&stockdata).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// **BUAT STOCK BARU JIKA BELUM ADA**
			if operation == "ADD" {
				newStock := stock.Stock{
					ProductID:    detail.ProductID,
					Quantity:     detail.Quantity,
					ExpiryDate:   detail.ExpiryDate,
					MinimumStock: 10, // Atur sesuai kebutuhan
				}
				return tx.Create(&newStock).Error
			}
			// Jika operation SUBTRACT tapi stock tidak ada, return error
			return fmt.Errorf("stock not found for product %s batch %s", detail.ProductCode, detail.BatchNumber)
		}
		return err
	}

	// **UPDATE STOCK YANG SUDAH ADA**
	var newQuantity int
	switch operation {
	case "ADD":
		newQuantity = stockdata.Quantity + detail.Quantity
	case "SUBTRACT":
		newQuantity = stockdata.Quantity - detail.Quantity
		if newQuantity < 0 {
			return fmt.Errorf("insufficient stock for product %s batch %s", detail.ProductCode, detail.BatchNumber)
		}
	}

	return tx.Model(&stockdata).Updates(map[string]interface{}{
		"Quantity":   newQuantity,
		"ExpiryDate": detail.ExpiryDate,
	}).Error
}
