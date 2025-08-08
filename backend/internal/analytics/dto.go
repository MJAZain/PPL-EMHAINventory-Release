package analytics

import "time"

// TimeRange represents the time range options
type TimeRange string

const (
	TimeRangeWeekly  TimeRange = "weekly"
	TimeRangeMonthly TimeRange = "monthly"
	TimeRangeYearly  TimeRange = "yearly"
)

// SalesAnalyticsRequest represents the request for sales analytics
type SalesAnalyticsRequest struct {
	TimeRange TimeRange  `json:"time_range" binding:"required,oneof=weekly monthly yearly"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// LineChartData represents data for line chart (sales trend)
type LineChartData struct {
	Date   string  `json:"date"`
	Total  float64 `json:"total"`
	Period string  `json:"period"`
}

// BarChartData represents data for bar chart (sales by category)
type BarChartData struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// ProductSalesData represents product sales data
type ProductSalesData struct {
	ProductID   uint    `json:"product_id"`
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	TotalQty    int     `json:"total_qty"`
	TotalAmount float64 `json:"total_amount"`
	SalesCount  int     `json:"sales_count"`
}

// SalesAnalyticsResponse represents the complete analytics response
type SalesAnalyticsResponse struct {
	LineChart     []LineChartData    `json:"line_chart"`
	BarChart      []BarChartData     `json:"bar_chart"`
	TopProducts   []ProductSalesData `json:"top_products"`
	LeastProducts []ProductSalesData `json:"least_products"`
	Summary       SalesSummary       `json:"summary"`
}

// SalesSummary represents sales summary data
type SalesSummary struct {
	TotalSales            float64 `json:"total_sales"`
	TotalTransactions     int     `json:"total_transactions"`
	AveragePerTransaction float64 `json:"average_per_transaction"`
	Period                string  `json:"period"`
}
