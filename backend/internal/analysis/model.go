package analysis

import "time"

type AnalysisRequest struct {
	Period    string `form:"period"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type AnalysisResponse struct {
	RequestParams      AnalysisRequest   `json:"request_params"`
	PeriodDescription  string            `json:"period_description"`
	StartDate          time.Time         `json:"start_date"`
	EndDate            time.Time         `json:"end_date"`
	TotalGrossRevenue  float64           `json:"total_gross_revenue"`
	TotalExpense       float64           `json:"total_expense"`
	NetProfit          float64           `json:"net_profit"`
	ProfitLossCompare  ProfitLossCompare `json:"profit_loss_compare"`
	PieChart           PieChartData      `json:"pie_chart"`
	BarChartRevenue    BarChartData      `json:"bar_chart_revenue"`
	BarChartExpense    BarChartData      `json:"bar_chart_expense"`
	RevenueTimeline    TimelineData      `json:"revenue_timeline"`
	ExpenseTimeline    TimelineData      `json:"expense_timeline"`
	TopSellingProducts []TopProduct      `json:"top_selling_products"`
}

type ProfitLossCompare struct {
	Status            string  `json:"status"`
	Difference        float64 `json:"difference"`
	Percentage        float64 `json:"percentage"`
	CurrentNetProfit  float64 `json:"current_net_profit"`
	PreviousNetProfit float64 `json:"previous_net_profit"`
	Message           string  `json:"message"`
}

type PieChartData struct {
	Labels []string  `json:"labels"`
	Values []float64 `json:"values"`
}

type BarChartData struct {
	Labels []string  `json:"labels"`
	Values []float64 `json:"values"`
}

type TimelineDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type TimelineData struct {
	Interval string              `json:"interval"`
	Labels   []string            `json:"labels"`
	Values   []float64           `json:"values"`
	Data     []TimelineDataPoint `json:"data"`
}

type TopProduct struct {
	ProductID     uint    `json:"product_id"`
	ProductName   string  `json:"product_name"`
	ProductCode   string  `json:"product_code"`
	TotalQuantity int     `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type revenueDetail struct {
	Source string
	Total  float64
}

type expenseDetail struct {
	Category string
	Total    float64
}

type timelineQueryResult struct {
	Date  string
	Value float64
}
