package controller

import (
	"fmt"
	"go-gin-auth/dto"
	"go-gin-auth/internal/opname"
	"go-gin-auth/mapper"
	"go-gin-auth/service"
	"go-gin-auth/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StockOpnameController struct {
	service service.StockOpnameService
}

func NewStockOpnameController(s service.StockOpnameService) *StockOpnameController {
	return &StockOpnameController{s}
}

func (c *StockOpnameController) Create(ctx *gin.Context) {

	var input dto.StockOpnameRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// Pemetaan dto ke model menggunakan mapper
	stockOpname := mapper.ToModelStockOpname(input)

	if err := c.service.Create(&stockOpname); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.Respond(ctx, http.StatusCreated, "Success", nil, input)
}

func (c *StockOpnameController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, data)
}

func (c *StockOpnameController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.service.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, data)
}

func (c *StockOpnameController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Cek apakah data ada terlebih dahulu
	exists, err := c.service.IsExist(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	if !exists {
		utils.Respond(ctx, http.StatusNotFound, "error", "Data tidak ditemukan", nil)
		return
	}

	err = c.service.Delete(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, "Data berhasil dihapus")
}

// GetStockOpnameHistory godoc
// @Summary Get stock opname history
// @Description Get a list of all stock adjustments history
// @Tags stock-opname
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /stock-opname/history [get]
func (c *StockOpnameController) GetStockOpnameHistory(ctx *gin.Context) {

	history, err := c.service.GetStockOpnameHistory(ctx)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Failed to get stock opname history", err.Error(), nil)
		return
	}

	utils.Respond(ctx, http.StatusOK, "Stock opname history retrieved successfully", nil, history)
}

// GetStockDiscrepancies godoc
// @Summary Get stock discrepancies
// @Description Get a list of products with significant stock discrepancies
// @Tags stock-opname
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /stock-opname/discrepancies [get]
func (c *StockOpnameController) GetStockDiscrepancies(ctx *gin.Context) {
	discrepancies, err := c.service.GetStockDiscrepancies(ctx)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Failed to get stock discrepancies", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Stock discrepancies retrieved successfully", nil, discrepancies)
}

// AdjustProductStock godoc
// @Summary Adjust product stock
// @Description Adjust stock for a specific product based on actual physical count
// @Tags stock-opname
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param adjustment body models.StockAdjustmentRequest true "Stock adjustment details"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /stock-opname/products/{product_id} [put]
func (c *StockOpnameController) AdjustProductStock(ctx *gin.Context) {
	productID := ctx.Param("product_id")

	var req dto.StockAdjustmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "Invalid request format", err.Error(), nil)
		return
	}

	adjustment, err := c.service.AdjustProductStock(ctx, productID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		}
		utils.Respond(ctx, status, "Failed to adjust product stock", err.Error(), nil)
		return
	}

	// Tambahkan informasi produk yang ditambahkan ke respons
	responseData := map[string]interface{}{
		"product_id":      adjustment.ProductID,
		"previous_stock":  adjustment.PreviousStock,
		"actual_stock":    adjustment.AdjustedStock,
		"discrepancy":     adjustment.AdjustmentQuantity,
		"adjustment_note": adjustment.AdjustmentNote,
		"opname_date":     adjustment.AdjustmentDate,
		"performed_by":    adjustment.PerformedBy,
	}

	utils.Respond(ctx, http.StatusOK, "Product stock adjusted successfully", nil, responseData)
}

// CreateDraft creates a new stock opname draft
func (h *StockOpnameController) CreateDraft(ctx *gin.Context) {
	var req dto.CreateDraftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "Invalid request payload", err.Error(), nil)
		return
	}
	fmt.Println("Request payload:", utils.GetCurrentUserID(ctx))
	opname, err := h.service.CreateDraft(strconv.FormatUint(uint64(utils.GetCurrentUserID(ctx)), 10), req.Tanggal, req.Catatan)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "successfully", nil, opname)
}

// AddProductToDraft adds a product to a draft stock opname
func (h *StockOpnameController) AddProductToDraft(ctx *gin.Context) {
	opnameID := ctx.Param("opnameID")

	var req dto.AddProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "Invalid request payload", err.Error(), nil)
		return
	}

	detail, err := h.service.AddProductToDraft(opnameID, req.ProductID)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	response := map[string]interface{}{
		"opname_id":  detail.OpnameID,
		"detail_id":  detail.DetailID,
		"product_id": detail.ProductID,
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, response)
}

// StartOpname starts the stock opname process
func (h *StockOpnameController) StartOpname(ctx *gin.Context) {
	opnameID := ctx.Param("opnameID")

	opname, err := h.service.StartOpname(opnameID, strconv.FormatUint(uint64(utils.GetCurrentUserID(ctx)), 10))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	response := map[string]interface{}{
		"opname_id":   opname.OpnameID,
		"opname_date": opname.OpnameDate,
		"start_time":  opname.StartTime,
		"end_time":    opname.EndTime,
		"status":      opname.Status,
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, response)
}

// CompleteOpname completes the stock opname process
func (h *StockOpnameController) CompleteOpname(ctx *gin.Context) {
	opnameID := ctx.Param("opnameID")

	opname, err := h.service.CompleteOpname(opnameID, strconv.FormatUint(uint64(utils.GetCurrentUserID(ctx)), 10))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	response := map[string]interface{}{
		"opname_id":   opname.OpnameID,
		"opname_date": opname.OpnameDate,
		"start_time":  opname.StartTime,
		"end_time":    opname.EndTime,
		"status":      opname.Status,
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, response)
}

// RecordActualStock records the actual stock count for a product
func (h *StockOpnameController) RecordActualStock(ctx *gin.Context) {
	detailIDStr := ctx.Param("detailID")

	detailID, err := strconv.Atoi(detailIDStr)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "Invalid detail ID", err.Error(), nil)
		return
	}

	var req dto.RecordStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "Invalid request payload", err.Error(), nil)
		return
	}

	detail, err := h.service.RecordActualStock(detailID, req.ActualStock, strconv.FormatUint(uint64(utils.GetCurrentUserID(ctx)), 10), req.Note)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	response := map[string]interface{}{
		"detail_id":              detail.DetailID,
		"opname_id":              detail.OpnameID,
		"product_id":             detail.ProductID,
		"system_stock":           detail.SystemStock,
		"actual_stock":           detail.ActualStock,
		"discrepancy":            detail.Discrepancy,
		"discrepancy_percentage": detail.DiscrepancyPercentage,
		"adjustment_note":        detail.AdjustmentNote,
	}
	utils.Respond(ctx, http.StatusOK, "Success", nil, response)
}

// GetDraft retrieves a draft stock opname
func (h *StockOpnameController) GetDraft(c *gin.Context) {
	opnameID := c.Param("opnameID")

	opname, err := h.service.GetDraft(opnameID)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, opname)
}

// UpdateDraft updates a draft stock opname
func (h *StockOpnameController) UpdateDraft(c *gin.Context) {
	opnameID := c.Param("opnameID")

	var req dto.UpdateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request payload", err.Error(), nil)
		return
	}

	opname, err := h.service.UpdateDraft(opnameID, req.Tanggal, req.Catatan)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Suksess", nil, opname)
}

// DeleteDraft deletes a draft stock opname
func (h *StockOpnameController) DeleteDraft(c *gin.Context) {
	opnameID := c.Param("opnameID")

	if err := h.service.DeleteDraft(opnameID); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Suksess", nil, nil)
}

// RemoveProductFromDraft removes a product from a draft stock opname
func (h *StockOpnameController) RemoveProductFromDraft(c *gin.Context) {
	opnameID := c.Param("opnameID")
	detailIDStr := c.Param("detailID")

	detailID, err := strconv.Atoi(detailIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid detail ID", err.Error(), nil)
		return
	}

	if err := h.service.RemoveProductFromDraft(opnameID, detailID); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Suksess", nil, nil)
}

// CancelOpname cancels the stock opname process
func (h *StockOpnameController) CancelOpname(c *gin.Context) {
	opnameID := c.Param("opnameID")

	opname, err := h.service.CancelOpname(opnameID, strconv.FormatUint(uint64(utils.GetCurrentUserID(c)), 10))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Suksess", nil, opname)
}

// GetOpnameDetails gets detailed information for a stock opname
func (h *StockOpnameController) GetOpnameDetails(c *gin.Context) {
	opnameID := c.Param("opnameID")

	opname, err := h.service.GetOpnameDetails(opnameID)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Suksess", nil, opname)
}

// GetOpnameList gets a list of stock opnames by status and date range
func (h *StockOpnameController) GetOpnameList(c *gin.Context) {
	opnameIdParam := c.Query("opname_id")
	statusParam := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var status opname.StockOpnameStatus
	if statusParam != "" {
		status = opname.StockOpnameStatus(statusParam)
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.Respond(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD", err.Error(), nil)
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			utils.Respond(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD", err.Error(), nil)
			return
		}
		// Set end date to end of day
		endDate = endDate.Add(24*time.Hour - time.Second)
	}

	opnames, err := h.service.GetOpnameList(status, startDate, endDate, opnameIdParam)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Suksess", nil, opnames)
}

func (c *StockOpnameController) GetProducts(ctx *gin.Context) {
	products, err := c.service.GetProducts(ctx.Request.Context())
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Suksess", nil, products)
}
