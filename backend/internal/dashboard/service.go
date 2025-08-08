// package dashboard

// import (
// 	"go-gin-auth/internal/prescription"
// 	"go-gin-auth/internal/sales"
// 	"time"

// 	"gorm.io/gorm"
// )

// type DashboardService struct {
// 	db *gorm.DB
// }

// func NewDashboardService(db *gorm.DB) *DashboardService {
// 	return &DashboardService{db: db}
// }

// // DashboardData represents the complete dashboard data
// type DashboardData struct {
// 	TotalSalesRegular      float64 `json:"total_sales_regular"`
// 	TotalSalesPrescription float64 `json:"total_sales_prescription"`
// 	TotalRevenue           float64 `json:"total_revenue"`
// 	Date                   string  `json:"date"`
// }

// // GetTotalSalesRegular gets total sales without prescription for today
// func (s *DashboardService) GetTotalSalesRegular(date time.Time) (float64, error) {
// 	var total float64

// 	// Get start and end of the day
// 	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
// 	endOfDay := startOfDay.Add(24 * time.Hour)

// 	err := s.db.Model(&sales.SalesRegular{}).
// 		Select("COALESCE(SUM(total_pay), 0)").
// 		Where("transaction_date >= ? AND transaction_date < ? AND deleted_at IS NULL", startOfDay, endOfDay).
// 		Scan(&total).Error

// 	return total, err
// }

// // GetTotalSalesPrescription gets total sales with prescription for today
// func (s *DashboardService) GetTotalSalesPrescription(date time.Time) (float64, error) {
// 	var total float64

// 	// Get start and end of the day
// 	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
// 	endOfDay := startOfDay.Add(24 * time.Hour)

// 	err := s.db.Model(&prescription.PrescriptionSale{}).
// 		Select("COALESCE(SUM(total_amount), 0)").
// 		Where("transaction_date >= ? AND transaction_date < ? AND deleted_at IS NULL", startOfDay, endOfDay).
// 		Scan(&total).Error

// 	return total, err
// }

// // GetTotalRevenue gets total revenue (prescription + regular sales - expenses)
// func (s *DashboardService) GetTotalRevenue(date time.Time) (float64, error) {
// 	regularTotal, err := s.GetTotalSalesRegular(date)
// 	if err != nil {
// 		return 0, err
// 	}

// 	prescriptionTotal, err := s.GetTotalSalesPrescription(date)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Get total expenses for the day (assuming you have an expenses table)
// 	totalExpenses, err := s.getTotalExpenses(date)
// 	if err != nil {
// 		return 0, err
// 	}

// 	totalRevenue := regularTotal + prescriptionTotal - totalExpenses
// 	return totalRevenue, nil
// }

// // getTotalExpenses gets total expenses for the day
// // Note: You need to implement this based on your expenses table structure
// func (s *DashboardService) getTotalExpenses(date time.Time) (float64, error) {
// 	var total float64

// 	// Get start and end of the day
// 	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
// 	endOfDay := startOfDay.Add(24 * time.Hour)

// 	// Assuming you have an expenses table - adjust the table name and columns as needed
// 	err := s.db.Table("expenses").
// 		Select("COALESCE(SUM(amount), 0)").
// 		Where("expense_date >= ? AND expense_date < ? AND deleted_at IS NULL", startOfDay, endOfDay).
// 		Scan(&total).Error

// 	return total, err
// }

// // GetDashboardData gets all dashboard data for a specific date
// func (s *DashboardService) GetDashboardData(date time.Time) (*DashboardData, error) {
// 	regularTotal, err := s.GetTotalSalesRegular(date)
// 	if err != nil {
// 		return nil, err
// 	}

// 	prescriptionTotal, err := s.GetTotalSalesPrescription(date)
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalRevenue, err := s.GetTotalRevenue(date)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DashboardData{
// 		TotalSalesRegular:      regularTotal,
// 		TotalSalesPrescription: prescriptionTotal,
// 		TotalRevenue:           totalRevenue,
// 		Date:                   date.Format("2006-01-02"),
// 	}, nil
// }

package dashboard

import (
	"go-gin-auth/internal/prescription"
	"go-gin-auth/internal/sales"
	"time"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// DashboardData represents the complete dashboard data
type DashboardData struct {
	TotalSalesRegular      float64 `json:"total_sales_regular"`
	TotalSalesPrescription float64 `json:"total_sales_prescription"`
	TotalRevenue           float64 `json:"total_revenue"`
	Date                   string  `json:"date"`
}

// GetTotalSalesRegular gets total sales without prescription for today
func (s *DashboardService) GetTotalSalesRegular(date time.Time) (float64, error) {
	var total float64

	// Get start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := s.db.Model(&sales.SalesRegular{}).
		Select("COALESCE(SUM(total_pay), 0)").
		Where("transaction_date >= ? AND transaction_date < ? AND deleted_at IS NULL", startOfDay, endOfDay).
		Scan(&total).Error

	return total, err
}

// GetTotalSalesPrescription gets total sales with prescription for today
func (s *DashboardService) GetTotalSalesPrescription(date time.Time) (float64, error) {
	var total float64

	// Get start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := s.db.Model(&prescription.PrescriptionSale{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("transaction_date >= ? AND transaction_date < ? AND deleted_at IS NULL", startOfDay, endOfDay).
		Scan(&total).Error

	return total, err
}

// GetTotalRevenue gets total revenue (prescription + regular sales only)
func (s *DashboardService) GetTotalRevenue(date time.Time) (float64, error) {
	regularTotal, err := s.GetTotalSalesRegular(date)
	if err != nil {
		return 0, err
	}

	prescriptionTotal, err := s.GetTotalSalesPrescription(date)
	if err != nil {
		return 0, err
	}

	// Total revenue is just the sum of both sales types (no expenses)
	totalRevenue := regularTotal + prescriptionTotal
	return totalRevenue, nil
}

// GetDashboardData gets all dashboard data for a specific date
func (s *DashboardService) GetDashboardData(date time.Time) (*DashboardData, error) {
	regularTotal, err := s.GetTotalSalesRegular(date)
	if err != nil {
		return nil, err
	}

	prescriptionTotal, err := s.GetTotalSalesPrescription(date)
	if err != nil {
		return nil, err
	}

	totalRevenue, err := s.GetTotalRevenue(date)
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		TotalSalesRegular:      regularTotal,
		TotalSalesPrescription: prescriptionTotal,
		TotalRevenue:           totalRevenue,
		Date:                   date.Format("2006-01-02"),
	}, nil
}