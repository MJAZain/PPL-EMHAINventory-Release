package analysis

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

const (
	isoDateFormat    = "2006-01-02"
	TopProductsLimit = 10
)

type Service interface {
	GetAnalysis(req AnalysisRequest) (*AnalysisResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAnalysis(req AnalysisRequest) (*AnalysisResponse, error) {
	currentStart, currentEnd, periodDesc, err := s.resolveDateRange(req)
	if err != nil {
		return nil, err
	}

	currentMetrics, err := s.calculateMetricsForPeriod(currentStart, currentEnd)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung metrik periode saat ini: %w", err)
	}

	prevStart, prevEnd, hasPrevPeriod := s.getPreviousDateRange(currentStart, req)
	var prevMetrics *metrics
	if hasPrevPeriod {
		prevMetrics, err = s.calculateMetricsForPeriod(prevStart, prevEnd)
		if err != nil {
			log.Printf("Peringatan: Gagal menghitung metrik periode sebelumnya: %v", err)
		}
	}

	durationDays := currentEnd.Sub(currentStart).Hours() / 24
	timelineInterval := "daily"
	if durationDays > 45 {
		timelineInterval = "monthly"
	}
	revenueTimeline, _ := s.repo.GetRevenueTimeline(currentStart, currentEnd, timelineInterval)
	expenseTimeline, _ := s.repo.GetExpenseTimeline(currentStart, currentEnd, timelineInterval)

	response := s.buildResponse(req, periodDesc, currentStart, currentEnd, currentMetrics, prevMetrics, revenueTimeline, expenseTimeline)
	response.RevenueTimeline = s.createTimelineData(revenueTimeline, timelineInterval)
	response.ExpenseTimeline = s.createTimelineData(expenseTimeline, timelineInterval)

	return response, nil
}

type metrics struct {
	TotalGrossRevenue float64
	TotalExpense      float64
	NetProfit         float64
	RevenueDetails    []revenueDetail
	ExpenseDetails    []expenseDetail
	TopProducts       []TopProduct
}

func (s *service) calculateMetricsForPeriod(start, end time.Time) (*metrics, error) {
	revenueDetails, err := s.repo.GetTotalRevenue(start, end)
	if err != nil {
		return nil, err
	}
	expenseDetails, err := s.repo.GetExpenseBreakdown(start, end)
	if err != nil {
		return nil, err
	}
	topProducts, err := s.repo.GetTopSellingProducts(start, end, TopProductsLimit)
	if err != nil {
		return nil, err
	}

	m := &metrics{
		RevenueDetails: revenueDetails,
		ExpenseDetails: expenseDetails,
		TopProducts:    topProducts,
	}

	for _, detail := range revenueDetails {
		m.TotalGrossRevenue += detail.Total
	}
	for _, detail := range expenseDetails {
		m.TotalExpense += detail.Total
	}
	m.NetProfit = m.TotalGrossRevenue - m.TotalExpense

	return m, nil
}

func (s *service) buildResponse(req AnalysisRequest, periodDesc string, start, end time.Time, current *metrics, previous *metrics, revTimeline, expTimeline []timelineQueryResult) *AnalysisResponse {
	return &AnalysisResponse{
		RequestParams:      req,
		PeriodDescription:  periodDesc,
		StartDate:          start,
		EndDate:            end,
		TotalGrossRevenue:  current.TotalGrossRevenue,
		TotalExpense:       current.TotalExpense,
		NetProfit:          current.NetProfit,
		ProfitLossCompare:  s.compareProfitLoss(current.NetProfit, previous),
		PieChart:           s.createPieChartData(current.TotalGrossRevenue, current.TotalExpense),
		BarChartRevenue:    s.createRevenueBarChartData(current.RevenueDetails),
		BarChartExpense:    s.createExpenseBarChartData(current.ExpenseDetails),
		TopSellingProducts: current.TopProducts,
	}
}

func (s *service) compareProfitLoss(currentProfit float64, prevMetrics *metrics) ProfitLossCompare {
	if prevMetrics == nil {
		return ProfitLossCompare{Status: "N/A", Message: "Tidak ada data periode sebelumnya untuk perbandingan."}
	}
	previousProfit := prevMetrics.NetProfit
	diff := currentProfit - previousProfit

	plc := ProfitLossCompare{
		CurrentNetProfit:  currentProfit,
		PreviousNetProfit: previousProfit,
		Difference:        diff,
	}

	if diff > 0 {
		plc.Status = "Naik"
	} else if diff < 0 {
		plc.Status = "Turun"
	} else {
		plc.Status = "Sama"
	}

	if previousProfit != 0 {
		plc.Percentage = (diff / math.Abs(previousProfit)) * 100
	} else if currentProfit > 0 {
		plc.Percentage = 100.0
		plc.Message = "Pendapatan tumbuh dari nol."
	} else {
		plc.Percentage = 0.0
	}
	return plc
}

func (s *service) createPieChartData(revenue, expense float64) PieChartData {
	return PieChartData{
		Labels: []string{"Pendapatan Kotor", "Pengeluaran"},
		Values: []float64{revenue, expense},
	}
}

func (s *service) createRevenueBarChartData(details []revenueDetail) BarChartData {
	var labels []string
	var values []float64
	for _, detail := range details {
		labels = append(labels, detail.Source)
		values = append(values, detail.Total)
	}
	return BarChartData{Labels: labels, Values: values}
}

func (s *service) createExpenseBarChartData(details []expenseDetail) BarChartData {
	var labels []string
	var values []float64
	for _, detail := range details {
		labels = append(labels, detail.Category)
		values = append(values, detail.Total)
	}
	return BarChartData{Labels: labels, Values: values}
}

func (s *service) createTimelineData(data []timelineQueryResult, interval string) TimelineData {
	td := TimelineData{
		Interval: interval,
		Data:     make([]TimelineDataPoint, len(data)),
		Labels:   make([]string, len(data)),
		Values:   make([]float64, len(data)),
	}
	for i, item := range data {
		td.Data[i] = TimelineDataPoint{Date: item.Date, Value: item.Value}
		td.Labels[i] = item.Date
		td.Values[i] = item.Value
	}
	return td
}

func (s *service) resolveDateRange(req AnalysisRequest) (start, end time.Time, desc string, err error) {
	now := time.Now()
	loc := now.Location()

	if req.StartDate != "" && req.EndDate != "" {
		start, err = time.ParseInLocation(isoDateFormat, req.StartDate, loc)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("format start_date tidak valid, gunakan YYYY-MM-DD")
		}
		end, err = time.ParseInLocation(isoDateFormat, req.EndDate, loc)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("format end_date tidak valid, gunakan YYYY-MM-DD")
		}
		end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		desc = fmt.Sprintf("Periode Kustom (%s - %s)", req.StartDate, req.EndDate)
		return
	}

	period := req.Period
	if period == "" {
		period = "monthly"
	}

	switch period {
	case "weekly":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -weekday+1)
		end = now
		desc = "Minggu Ini"
	case "monthly":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		end = now
		desc = "Bulan Ini"
	case "yearly":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		end = now
		desc = "Tahun Ini"
	default:
		err = fmt.Errorf("periode tidak valid: '%s'. Pilih antara 'weekly', 'monthly', 'yearly', atau gunakan start_date & end_date", period)
	}
	return
}

func (s *service) getPreviousDateRange(currentStart time.Time, req AnalysisRequest) (start, end time.Time, ok bool) {
	if req.StartDate != "" && req.EndDate != "" {
		return time.Time{}, time.Time{}, false
	}

	ok = true
	period := req.Period
	if period == "" {
		period = "monthly"
	}

	switch period {
	case "weekly":
		end = currentStart.Add(-time.Second)
		start = currentStart.AddDate(0, 0, -7)
	case "monthly":
		end = currentStart.Add(-time.Second)
		start = currentStart.AddDate(0, -1, 0)
	case "yearly":
		end = currentStart.Add(-time.Second)
		start = currentStart.AddDate(-1, 0, 0)
	default:
		ok = false
	}
	return
}
