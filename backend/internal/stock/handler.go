package stock

import (
	"go-gin-auth/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	Service *StockService
}

func NewStockHandler(s *StockService) *StockHandler {
	return &StockHandler{Service: s}
}

func (h *StockHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/current", h.GetCurrent)
	r.GET("/batches", h.GetBatches)
	r.GET("/low", h.GetLowStock)
	r.GET("/expiring-soon", h.GetExpiringSoon)
	r.GET("/summary", h.GetSummary)
	r.GET("/:item_id", h.GetDetail)
}

func (h *StockHandler) GetCurrent(c *gin.Context) {
	data, err := h.Service.GetCurrentStocks()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, data)
}

func (h *StockHandler) GetBatches(c *gin.Context) {
	var id *uint
	if itemID := c.Query("item_id"); itemID != "" {
		val, err := strconv.Atoi(itemID)
		if err == nil {
			u := uint(val)
			id = &u
		}
	}
	data, err := h.Service.GetStockBatches(id)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, data)
}

func (h *StockHandler) GetLowStock(c *gin.Context) {
	data, err := h.Service.GetLowStock()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, data)
}

func (h *StockHandler) GetExpiringSoon(c *gin.Context) {
	monthsStr := c.Query("months")
	months := 3 // default

	if monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 {
			months = m
		}
	}

	data, err := h.Service.GetExpiringSoonStocks(months)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Success", nil, data)
}

func (h *StockHandler) GetSummary(c *gin.Context) {
	data, err := h.Service.GetStockSummary()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, data)
}

func (h *StockHandler) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}
	data, err := h.Service.GetStockDetail(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Success", nil, data)
}
