// package dashboard

// import (
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// type DashboardHandler struct {
// 	dashboardService *DashboardService
// }

// func NewDashboardHandler(dashboardService *DashboardService) *DashboardHandler {
// 	return &DashboardHandler{
// 		dashboardService: dashboardService,
// 	}
// }

// // Response structures
// type APIResponse struct {
// 	Status  bool        `json:"status"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// type SalesRegularResponse struct {
// 	Total float64 `json:"total"`
// 	Date  string  `json:"date"`
// }

// type SalesPrescriptionResponse struct {
// 	Total float64 `json:"total"`
// 	Date  string  `json:"date"`
// }

// type RevenueResponse struct {
// 	TotalRevenue           float64 `json:"total_revenue"`
// 	TotalSalesRegular      float64 `json:"total_sales_regular"`
// 	TotalSalesPrescription float64 `json:"total_sales_prescription"`
// 	TotalExpenses          float64 `json:"total_expenses"`
// 	Date                   string  `json:"date"`
// }

// // GetTotalSalesRegular gets total sales without prescription for today
// // @Summary Get total sales regular
// // @Description Get total sales without prescription for a specific date (default: today)
// // @Tags Dashboard
// // @Accept json
// // @Produce json
// // @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// // @Success 200 {object} APIResponse{data=SalesRegularResponse}
// // @Failure 400 {object} APIResponse
// // @Failure 500 {object} APIResponse
// // @Router /api/dashboard/sales-regular [get]
// func (h *DashboardHandler) GetTotalSalesRegular(c *gin.Context) {
// 	// Parse date from query parameter, default to today
// 	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

// 	date, err := time.Parse("2006-01-02", dateStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, APIResponse{
// 			Status:  false,
// 			Message: "Invalid date format. Use YYYY-MM-DD",
// 		})
// 		return
// 	}

// 	total, err := h.dashboardService.GetTotalSalesRegular(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to get total sales regular",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, APIResponse{
// 		Status:  true,
// 		Message: "Total sales regular retrieved successfully",
// 		Data: SalesRegularResponse{
// 			Total: total,
// 			Date:  dateStr,
// 		},
// 	})
// }

// // GetTotalSalesPrescription gets total sales with prescription for today
// // @Summary Get total sales prescription
// // @Description Get total sales with prescription for a specific date (default: today)
// // @Tags Dashboard
// // @Accept json
// // @Produce json
// // @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// // @Success 200 {object} APIResponse{data=SalesPrescriptionResponse}
// // @Failure 400 {object} APIResponse
// // @Failure 500 {object} APIResponse
// // @Router /api/dashboard/sales-prescription [get]
// func (h *DashboardHandler) GetTotalSalesPrescription(c *gin.Context) {
// 	// Parse date from query parameter, default to today
// 	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

// 	date, err := time.Parse("2006-01-02", dateStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, APIResponse{
// 			Status:  false,
// 			Message: "Invalid date format. Use YYYY-MM-DD",
// 		})
// 		return
// 	}

// 	total, err := h.dashboardService.GetTotalSalesPrescription(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to get total sales prescription",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, APIResponse{
// 		Status:  true,
// 		Message: "Total sales prescription retrieved successfully",
// 		Data: SalesPrescriptionResponse{
// 			Total: total,
// 			Date:  dateStr,
// 		},
// 	})
// }

// // GetTotalRevenue gets total revenue (prescription + regular sales - expenses)
// // @Summary Get total revenue
// // @Description Get total revenue (prescription + regular sales - expenses) for a specific date (default: today)
// // @Tags Dashboard
// // @Accept json
// // @Produce json
// // @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// // @Success 200 {object} APIResponse{data=RevenueResponse}
// // @Failure 400 {object} APIResponse
// // @Failure 500 {object} APIResponse
// // @Router /api/dashboard/revenue [get]
// func (h *DashboardHandler) GetTotalRevenue(c *gin.Context) {
// 	// Parse date from query parameter, default to today
// 	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

// 	date, err := time.Parse("2006-01-02", dateStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, APIResponse{
// 			Status:  false,
// 			Message: "Invalid date format. Use YYYY-MM-DD",
// 		})
// 		return
// 	}

// 	// Get individual totals
// 	regularTotal, err := h.dashboardService.GetTotalSalesRegular(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to get regular sales total",
// 		})
// 		return
// 	}

// 	prescriptionTotal, err := h.dashboardService.GetTotalSalesPrescription(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to get prescription sales total",
// 		})
// 		return
// 	}

// 	totalRevenue, err := h.dashboardService.GetTotalRevenue(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to calculate total revenue",
// 		})
// 		return
// 	}

// 	// Calculate expenses (revenue - sales)
// 	totalExpenses := (regularTotal + prescriptionTotal) - totalRevenue

// 	c.JSON(http.StatusOK, APIResponse{
// 		Status:  true,
// 		Message: "Total revenue retrieved successfully",
// 		Data: RevenueResponse{
// 			TotalRevenue:           totalRevenue,
// 			TotalSalesRegular:      regularTotal,
// 			TotalSalesPrescription: prescriptionTotal,
// 			TotalExpenses:          totalExpenses,
// 			Date:                   dateStr,
// 		},
// 	})
// }

// // GetDashboardSummary gets all dashboard data in one request
// // @Summary Get dashboard summary
// // @Description Get complete dashboard data (regular sales, prescription sales, and revenue) for a specific date
// // @Tags Dashboard
// // @Accept json
// // @Produce json
// // @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// // @Success 200 {object} APIResponse{data=DashboardData}
// // @Failure 400 {object} APIResponse
// // @Failure 500 {object} APIResponse
// // @Router /api/dashboard/summary [get]
// func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
// 	// Parse date from query parameter, default to today
// 	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

// 	date, err := time.Parse("2006-01-02", dateStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, APIResponse{
// 			Status:  false,
// 			Message: "Invalid date format. Use YYYY-MM-DD",
// 		})
// 		return
// 	}

// 	dashboardData, err := h.dashboardService.GetDashboardData(date)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, APIResponse{
// 			Status:  false,
// 			Message: "Failed to get dashboard data",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, APIResponse{
// 		Status:  true,
// 		Message: "Dashboard data retrieved successfully",
// 		Data:    dashboardData,
// 	})
// }

// // SetupDashboardRoutes sets up the dashboard routes
//
//	func (h *DashboardHandler) DashboardRoutes(r *gin.RouterGroup) {
//		r.GET("/sales-regular", h.GetTotalSalesRegular)
//		r.GET("/sales-prescription", h.GetTotalSalesPrescription)
//		r.GET("/revenue", h.GetTotalRevenue)
//		r.GET("/summary", h.GetDashboardSummary)
//	}
package dashboard

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *DashboardService
}

func NewDashboardHandler(dashboardService *DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// Response structures
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SalesRegularResponse struct {
	Total float64 `json:"total"`
	Date  string  `json:"date"`
}

type SalesPrescriptionResponse struct {
	Total float64 `json:"total"`
	Date  string  `json:"date"`
}

type RevenueResponse struct {
	TotalRevenue           float64 `json:"total_revenue"`
	TotalSalesRegular      float64 `json:"total_sales_regular"`
	TotalSalesPrescription float64 `json:"total_sales_prescription"`
	Date                   string  `json:"date"`
}

// GetTotalSalesRegular gets total sales without prescription for today
// @Summary Get total sales regular
// @Description Get total sales without prescription for a specific date (default: today)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// @Success 200 {object} APIResponse{data=SalesRegularResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/dashboard/sales-regular [get]
func (h *DashboardHandler) GetTotalSalesRegular(c *gin.Context) {
	// Parse date from query parameter, default to today
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  false,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	total, err := h.dashboardService.GetTotalSalesRegular(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to get total sales regular",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  true,
		Message: "Total sales regular retrieved successfully",
		Data: SalesRegularResponse{
			Total: total,
			Date:  dateStr,
		},
	})
}

// GetTotalSalesPrescription gets total sales with prescription for today
// @Summary Get total sales prescription
// @Description Get total sales with prescription for a specific date (default: today)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// @Success 200 {object} APIResponse{data=SalesPrescriptionResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/dashboard/sales-prescription [get]
func (h *DashboardHandler) GetTotalSalesPrescription(c *gin.Context) {
	// Parse date from query parameter, default to today
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  false,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	total, err := h.dashboardService.GetTotalSalesPrescription(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to get total sales prescription",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  true,
		Message: "Total sales prescription retrieved successfully",
		Data: SalesPrescriptionResponse{
			Total: total,
			Date:  dateStr,
		},
	})
}

// GetTotalRevenue gets total revenue (prescription + regular sales)
// @Summary Get total revenue
// @Description Get total revenue (prescription + regular sales) for a specific date (default: today)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// @Success 200 {object} APIResponse{data=RevenueResponse}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/dashboard/revenue [get]
func (h *DashboardHandler) GetTotalRevenue(c *gin.Context) {
	// Parse date from query parameter, default to today
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  false,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	// Get individual totals
	regularTotal, err := h.dashboardService.GetTotalSalesRegular(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to get regular sales total",
		})
		return
	}

	prescriptionTotal, err := h.dashboardService.GetTotalSalesPrescription(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to get prescription sales total",
		})
		return
	}

	totalRevenue, err := h.dashboardService.GetTotalRevenue(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to calculate total revenue",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  true,
		Message: "Total revenue retrieved successfully",
		Data: RevenueResponse{
			TotalRevenue:           totalRevenue,
			TotalSalesRegular:      regularTotal,
			TotalSalesPrescription: prescriptionTotal,
			Date:                   dateStr,
		},
	})
}

// GetDashboardSummary gets all dashboard data in one request
// @Summary Get dashboard summary
// @Description Get complete dashboard data (regular sales, prescription sales, and revenue) for a specific date
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format (default: today)"
// @Success 200 {object} APIResponse{data=DashboardData}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/dashboard/summary [get]
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	// Parse date from query parameter, default to today
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  false,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
		return
	}

	dashboardData, err := h.dashboardService.GetDashboardData(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  false,
			Message: "Failed to get dashboard data",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  true,
		Message: "Dashboard data retrieved successfully",
		Data:    dashboardData,
	})
}

// DashboardRoutes sets up the dashboard routes
func (h *DashboardHandler) DashboardRoutes(r *gin.RouterGroup) {
	r.GET("/sales-regular", h.GetTotalSalesRegular)
	r.GET("/sales-prescription", h.GetTotalSalesPrescription)
	r.GET("/revenue", h.GetTotalRevenue)
	r.GET("/summary", h.GetDashboardSummary)
}
