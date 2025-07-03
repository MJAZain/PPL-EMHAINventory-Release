package pbf

import (
	"fmt"
	"go-gin-auth/config"
	"go-gin-auth/internal/product"
	"go-gin-auth/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	// Calculate total purchase
	var totalPurchase float64
	var details []IncomingPBFDetail

	for _, detailReq := range req.Details {
		// Get product info
		var product product.Product
		if err := config.DB.First(&product, detailReq.ProductID).Error; err != nil {
			utils.Respond(c, http.StatusBadRequest, fmt.Sprintf("Product with ID %d not found", detailReq.ProductID), err.Error(), nil)
			return
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
			Unit:          product.Unit.Name,
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

	// Check if record exists
	var existingRecord IncomingPBF
	if err := config.DB.First(&existingRecord, id).Error; err != nil {
		utils.Respond(c, http.StatusNotFound, "Record not found", err.Error(), nil)
		return
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
		config.DB.First(&product, detailReq.ProductID)

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
			Unit:          product.Unit.Name,
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
	for _, detail := range details {
		config.DB.Create(&detail)
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

	// Check if record exists
	var existingRecord IncomingPBF
	if err := config.DB.First(&existingRecord, id).Error; err != nil {
		utils.Respond(c, http.StatusNotFound, "Record not found", err.Error(), nil)
		return
	}

	// Delete details first (due to foreign key constraint)
	config.DB.Where("incoming_pbf_id = ?", id).Delete(&IncomingPBFDetail{})

	// Delete main record
	if err := config.DB.Delete(&existingRecord).Error; err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to delete record", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Record deleted successfully", nil, nil)
}
