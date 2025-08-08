package analytics

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

type SalesAnalyticsService struct {
	db *gorm.DB
}

func NewSalesAnalyticsService(db *gorm.DB) *SalesAnalyticsService {
	return &SalesAnalyticsService{db: db}
}

// GetSalesAnalytics returns complete sales analytics data
func (s *SalesAnalyticsService) GetSalesAnalytics(req SalesAnalyticsRequest) (*SalesAnalyticsResponse, error) {
	startDate, endDate := s.calculateDateRange(req.TimeRange, req.StartDate, req.EndDate)

	lineChart, err := s.getLineChartData(req.TimeRange, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get line chart data: %w", err)
	}

	barChart, err := s.getBarChartData(req.TimeRange, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get bar chart data: %w", err)
	}

	topProducts, err := s.getTopProducts(startDate, endDate, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}

	leastProducts, err := s.getLeastProducts(startDate, endDate, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get least products: %w", err)
	}

	summary, err := s.getSalesSummary(startDate, endDate, string(req.TimeRange))
	if err != nil {
		return nil, fmt.Errorf("failed to get sales summary: %w", err)
	}

	return &SalesAnalyticsResponse{
		LineChart:     lineChart,
		BarChart:      barChart,
		TopProducts:   topProducts,
		LeastProducts: leastProducts,
		Summary:       summary,
	}, nil
}

// getLineChartData returns line chart data for sales trend
func (s *SalesAnalyticsService) getLineChartData(timeRange TimeRange, startDate, endDate time.Time) ([]LineChartData, error) {
	var results []LineChartData
	var groupBy, selectDate string

	// Tentukan grouping dan format tanggal sesuai PostgreSQL
	switch timeRange {
	case TimeRangeWeekly:
		groupBy = "TO_CHAR(transaction_date, 'YYYY-MM-DD')"
		selectDate = "TO_CHAR(transaction_date, 'YYYY-MM-DD')"
	case TimeRangeMonthly:
		groupBy = "TO_CHAR(transaction_date, 'YYYY-MM')"
		selectDate = "TO_CHAR(transaction_date, 'YYYY-MM')"
	case TimeRangeYearly:
		groupBy = "TO_CHAR(transaction_date, 'YYYY')"
		selectDate = "TO_CHAR(transaction_date, 'YYYY')"
	default:
		return nil, fmt.Errorf("invalid time range")
	}

	// Query untuk sales_regulars
	regularQuery := fmt.Sprintf(`
		SELECT 
			%s AS date,
			COALESCE(SUM(total_pay), 0) AS total,
			'regular' AS period
		FROM sales_regulars 
		WHERE transaction_date BETWEEN ? AND ?
		AND deleted_at IS NULL
		GROUP BY %s
		ORDER BY date ASC
	`, selectDate, groupBy)

	// Query untuk prescription_sales
	prescriptionQuery := fmt.Sprintf(`
		SELECT 
			%s AS date,
			COALESCE(SUM(total_amount), 0) AS total,
			'prescription' AS period
		FROM prescription_sales 
		WHERE transaction_date BETWEEN ? AND ?
		AND deleted_at IS NULL
		GROUP BY %s
		ORDER BY date ASC
	`, selectDate, groupBy)

	// Struct untuk menampung hasil query
	var regularResults []struct {
		Date   string  `json:"date"`
		Total  float64 `json:"total"`
		Period string  `json:"period"`
	}
	var prescriptionResults []struct {
		Date   string  `json:"date"`
		Total  float64 `json:"total"`
		Period string  `json:"period"`
	}

	// Eksekusi query sales_regulars
	if err := s.db.Raw(regularQuery, startDate, endDate).Scan(&regularResults).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch regular sales: %w", err)
	}

	// Eksekusi query prescription_sales
	if err := s.db.Raw(prescriptionQuery, startDate, endDate).Scan(&prescriptionResults).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch prescription sales: %w", err)
	}

	// Gabungkan hasil dari kedua sumber ke map berdasarkan tanggal
	dateMap := make(map[string]float64)
	for _, r := range regularResults {
		dateMap[r.Date] += r.Total
	}
	for _, r := range prescriptionResults {
		dateMap[r.Date] += r.Total
	}

	// Ubah map ke slice hasil akhir
	for date, total := range dateMap {
		results = append(results, LineChartData{
			Date:   date,
			Total:  total,
			Period: string(timeRange),
		})
	}

	// Optional: urutkan hasil berdasarkan tanggal ascending
	// Karena map tidak berurut, urutkan manual:
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date < results[j].Date
	})

	return results, nil
}

func (s *SalesAnalyticsService) getBarChartData(timeRange TimeRange, startDate, endDate time.Time) ([]BarChartData, error) {
	var results []BarChartData

	regularQuery := `
		SELECT 
			COALESCE(c.name, 'Uncategorized') as category,
			SUM(sri.sub_total) as total,
			COUNT(DISTINCT sr.id) as count
		FROM sales_regulars sr
		JOIN sales_regular_items sri ON sr.id = sri.sales_regular_id
		JOIN products p ON sri.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE sr.transaction_date BETWEEN ? AND ?
		AND sr.deleted_at IS NULL AND sri.deleted_at IS NULL AND p.deleted_at IS NULL
		GROUP BY c.name
	`

	// Query untuk prescription sales (kategori selalu 'Prescription' misal)
	prescriptionQuery := `
		SELECT 
			'Prescription' as category,
			SUM(pi.sub_total) as total,
			COUNT(DISTINCT ps.id) as count
		FROM prescription_sales ps
		JOIN prescription_items pi ON ps.id = pi.prescription_sale_id
		WHERE ps.transaction_date BETWEEN ? AND ?
		AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL
		GROUP BY category
	`

	var regularResults []BarChartData
	if err := s.db.Raw(regularQuery, startDate, endDate).Scan(&regularResults).Error; err != nil {
		return nil, err
	}

	var prescriptionResults []BarChartData
	if err := s.db.Raw(prescriptionQuery, startDate, endDate).Scan(&prescriptionResults).Error; err != nil {
		return nil, err
	}

	categoryMap := make(map[string]*BarChartData)

	for _, r := range regularResults {
		if exist, ok := categoryMap[r.Category]; ok {
			exist.Total += r.Total
			exist.Count += r.Count
		} else {
			categoryMap[r.Category] = &BarChartData{
				Category: r.Category,
				Total:    r.Total,
				Count:    r.Count,
			}
		}
	}

	for _, r := range prescriptionResults {
		if exist, ok := categoryMap[r.Category]; ok {
			exist.Total += r.Total
			exist.Count += r.Count
		} else {
			categoryMap[r.Category] = &BarChartData{
				Category: r.Category,
				Total:    r.Total,
				Count:    r.Count,
			}
		}
	}

	for _, v := range categoryMap {
		results = append(results, *v)
	}

	return results, nil
}

// getTopProducts returns top selling products
func (s *SalesAnalyticsService) getTopProducts(startDate, endDate time.Time, limit int) ([]ProductSalesData, error) {
	var results []ProductSalesData

	// Query for regular sales
	regularQuery := `
		SELECT 
			sri.product_id,
			sri.product_code,
			sri.product_name,
			SUM(sri.qty) as total_qty,
			SUM(sri.sub_total) as total_amount,
			COUNT(DISTINCT sr.id) as sales_count
		FROM sales_regulars sr
		JOIN sales_regular_items sri ON sr.id = sri.sales_regular_id
		WHERE sr.transaction_date >= ? AND sr.transaction_date <= ?
		AND sr.deleted_at IS NULL AND sri.deleted_at IS NULL
		GROUP BY sri.product_id, sri.product_code, sri.product_name
	`

	// Query for prescription sales
	prescriptionQuery := `
		SELECT 
			pi.stock_id as product_id,
			pi.item_code as product_code,
			pi.item_name as product_name,
			SUM(pi.quantity) as total_qty,
			SUM(pi.sub_total) as total_amount,
			COUNT(DISTINCT ps.id) as sales_count
		FROM prescription_sales ps
		JOIN prescription_items pi ON ps.id = pi.prescription_sale_id
		WHERE ps.transaction_date >= ? AND ps.transaction_date <= ?
		AND ps.deleted_at IS NULL AND pi.deleted_at IS NULL
		GROUP BY pi.stock_id, pi.item_code, pi.item_name
	`

	// Execute regular sales query
	var regularResults []ProductSalesData
	if err := s.db.Raw(regularQuery, startDate, endDate).Scan(&regularResults).Error; err != nil {
		return nil, err
	}

	// Execute prescription sales query
	var prescriptionResults []ProductSalesData
	if err := s.db.Raw(prescriptionQuery, startDate, endDate).Scan(&prescriptionResults).Error; err != nil {
		return nil, err
	}

	// Combine results
	productMap := make(map[uint]*ProductSalesData)

	for _, result := range regularResults {
		if existing, exists := productMap[result.ProductID]; exists {
			existing.TotalQty += result.TotalQty
			existing.TotalAmount += result.TotalAmount
			existing.SalesCount += result.SalesCount
		} else {
			productMap[result.ProductID] = &ProductSalesData{
				ProductID:   result.ProductID,
				ProductCode: result.ProductCode,
				ProductName: result.ProductName,
				TotalQty:    result.TotalQty,
				TotalAmount: result.TotalAmount,
				SalesCount:  result.SalesCount,
			}
		}
	}

	for _, result := range prescriptionResults {
		if existing, exists := productMap[result.ProductID]; exists {
			existing.TotalQty += result.TotalQty
			existing.TotalAmount += result.TotalAmount
			existing.SalesCount += result.SalesCount
		} else {
			productMap[result.ProductID] = &ProductSalesData{
				ProductID:   result.ProductID,
				ProductCode: result.ProductCode,
				ProductName: result.ProductName,
				TotalQty:    result.TotalQty,
				TotalAmount: result.TotalAmount,
				SalesCount:  result.SalesCount,
			}
		}
	}

	// Convert to slice and sort by total amount (descending)
	for _, data := range productMap {
		results = append(results, *data)
	}

	// Sort by total amount (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].TotalAmount < results[j].TotalAmount {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// getLeastProducts returns least selling products
func (s *SalesAnalyticsService) getLeastProducts(startDate, endDate time.Time, limit int) ([]ProductSalesData, error) {
	products, err := s.getTopProducts(startDate, endDate, 0) // Get all products
	if err != nil {
		return nil, err
	}

	// Sort by total amount (ascending)
	for i := 0; i < len(products)-1; i++ {
		for j := i + 1; j < len(products); j++ {
			if products[i].TotalAmount > products[j].TotalAmount {
				products[i], products[j] = products[j], products[i]
			}
		}
	}

	// Limit results
	if len(products) > limit {
		products = products[:limit]
	}

	return products, nil
}

// getSalesSummary returns sales summary data
func (s *SalesAnalyticsService) getSalesSummary(startDate, endDate time.Time, period string) (SalesSummary, error) {
	var summary SalesSummary

	// Query for regular sales summary
	regularQuery := `
		SELECT 
			COALESCE(SUM(total_pay), 0) as total_sales,
			COUNT(*) as total_transactions
		FROM sales_regulars 
		WHERE transaction_date >= ? AND transaction_date <= ?
		AND deleted_at IS NULL
	`

	// Query for prescription sales summary
	prescriptionQuery := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			COUNT(*) as total_transactions
		FROM prescription_sales 
		WHERE transaction_date >= ? AND transaction_date <= ?
		AND deleted_at IS NULL
	`

	var regularSummary struct {
		TotalSales        float64 `json:"total_sales"`
		TotalTransactions int     `json:"total_transactions"`
	}

	var prescriptionSummary struct {
		TotalSales        float64 `json:"total_sales"`
		TotalTransactions int     `json:"total_transactions"`
	}

	if err := s.db.Raw(regularQuery, startDate, endDate).Scan(&regularSummary).Error; err != nil {
		return summary, err
	}

	if err := s.db.Raw(prescriptionQuery, startDate, endDate).Scan(&prescriptionSummary).Error; err != nil {
		return summary, err
	}

	summary.TotalSales = regularSummary.TotalSales + prescriptionSummary.TotalSales
	summary.TotalTransactions = regularSummary.TotalTransactions + prescriptionSummary.TotalTransactions
	summary.Period = period

	if summary.TotalTransactions > 0 {
		summary.AveragePerTransaction = summary.TotalSales / float64(summary.TotalTransactions)
	}

	return summary, nil
}

// calculateDateRange calculates start and end dates based on time range
func (s *SalesAnalyticsService) calculateDateRange(timeRange TimeRange, startDate, endDate *time.Time) (time.Time, time.Time) {
	now := time.Now()

	if startDate != nil && endDate != nil {
		return *startDate, *endDate
	}

	switch timeRange {
	case TimeRangeWeekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		start := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
		end := start.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return start, end
	case TimeRangeMonthly:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return start, end
	case TimeRangeYearly:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(1, 0, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return start, end
	default:
		// Default to current month
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return start, end
	}
}
