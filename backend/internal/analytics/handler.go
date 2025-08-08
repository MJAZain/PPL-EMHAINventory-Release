package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SalesAnalyticsHandler struct {
	service *SalesAnalyticsService
}

func NewSalesAnalyticsHandler(service *SalesAnalyticsService) *SalesAnalyticsHandler {
	return &SalesAnalyticsHandler{service: service}
}

// GetSalesAnalytics handles the complete sales analytics request
// @Summary Get sales analytics
// @Description Get complete sales analytics including line chart, bar chart, top products, and least products
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} SalesAnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics [post]
func (h *SalesAnalyticsHandler) GetSalesAnalytics(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	result, err := h.service.GetSalesAnalytics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get sales analytics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetLineChartData handles line chart data request
// @Summary Get line chart data
// @Description Get sales trend data for line chart
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} []LineChartData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics/line-chart [post]
func (h *SalesAnalyticsHandler) GetLineChartData(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	startDate, endDate := h.service.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)
	result, err := h.service.getLineChartData(req.TimeRange, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get line chart data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetBarChartData handles bar chart data request
// @Summary Get bar chart data
// @Description Get sales data by category for bar chart
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} []BarChartData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics/bar-chart [post]
func (h *SalesAnalyticsHandler) GetBarChartData(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	startDate, endDate := h.service.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)
	result, err := h.service.getBarChartData(req.TimeRange, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get bar chart data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetTopProducts handles top products request
// @Summary Get top selling products
// @Description Get top 5 most selling products
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} []ProductSalesData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics/top-products [post]
func (h *SalesAnalyticsHandler) GetTopProducts(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	startDate, endDate := h.service.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)
	result, err := h.service.getTopProducts(startDate, endDate, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get top products",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetLeastProducts handles least products request
// @Summary Get least selling products
// @Description Get top 5 least selling products
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} []ProductSalesData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics/least-products [post]
func (h *SalesAnalyticsHandler) GetLeastProducts(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	startDate, endDate := h.service.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)
	result, err := h.service.getLeastProducts(startDate, endDate, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get least products",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetSalesSummary handles sales summary request
// @Summary Get sales summary
// @Description Get sales summary data
// @Tags Sales Analytics
// @Accept json
// @Produce json
// @Param request body SalesAnalyticsRequest true "Analytics request"
// @Success 200 {object} SalesSummary
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sales/analytics/summary [post]
func (h *SalesAnalyticsHandler) GetSalesSummary(c *gin.Context) {
	var req SalesAnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	startDate, endDate := h.service.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)
	result, err := h.service.getSalesSummary(startDate, endDate, string(req.TimeRange))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get sales summary",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// DashboardRoutes sets up the dashboard routes
func (h *SalesAnalyticsHandler) SetupAnalyticsRoutes(r *gin.RouterGroup) {
	// Individual endpoints
	r.POST("/line-chart", h.GetLineChartData)
	r.POST("/bar-chart", h.GetBarChartData)
	r.POST("/top-products", h.GetTopProducts)
	r.POST("/least-products", h.GetLeastProducts)
	r.POST("/summary", h.GetSalesSummary)
}
