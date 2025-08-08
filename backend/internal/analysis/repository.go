package analysis

import (
	"fmt"
	"go-gin-auth/internal/expense"
	"go-gin-auth/internal/prescription"
	"go-gin-auth/internal/sales"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	GetTotalRevenue(start, end time.Time) ([]revenueDetail, error)
	GetExpenseBreakdown(start, end time.Time) ([]expenseDetail, error)
	GetRevenueTimeline(start, end time.Time, interval string) ([]timelineQueryResult, error)
	GetExpenseTimeline(start, end time.Time, interval string) ([]timelineQueryResult, error)
	GetTopSellingProducts(start, end time.Time, limit int) ([]TopProduct, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a new instance of Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetTotalRevenue gets total revenue from prescription and regular sales
// for a given date range.
func (r *repository) GetTotalRevenue(start, end time.Time) ([]revenueDetail, error) {
	var results []revenueDetail
	var prescriptionTotal float64
	err := r.db.Model(&prescription.PrescriptionSale{}).
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total_amount), 0)").
		Row().Scan(&prescriptionTotal)
	if err != nil {
		return nil, fmt.Errorf("gagal query penjualan resep: %w", err)
	}
	results = append(results, revenueDetail{Source: "Penjualan Resep", Total: prescriptionTotal})

	var regularTotal float64
	err = r.db.Model(&sales.SalesRegular{}).
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total_pay), 0)").
		Row().Scan(&regularTotal)
	if err != nil {
		return nil, fmt.Errorf("gagal query penjualan reguler: %w", err)
	}
	results = append(results, revenueDetail{Source: "Penjualan Reguler", Total: regularTotal})

	return results, nil
}

// GetExpenseBreakdown gets the total expense for each expense type within a given date range.
// The result is a slice of expenseDetail, sorted in descending order of total expense.
func (r *repository) GetExpenseBreakdown(start, end time.Time) ([]expenseDetail, error) {
	var details []expenseDetail
	err := r.db.Model(&expense.Expense{}).
		Select("et.name as category, SUM(expenses.amount) as total").
		Joins("JOIN expense_types et ON et.id = expenses.expense_type_id").
		Where("expenses.date BETWEEN ? AND ?", start, end).
		Group("et.name").
		Order("total DESC").
		Scan(&details).Error
	return details, err
}

// getDateFormatString returns the SQL to extract the date part
// from a given field, given the database dialect and the desired
// interval. The interval can be "monthly" to get the year and month, or
// any other value to get the full date.
func getDateFormatString(db *gorm.DB, interval string) string {
	dialect := db.Dialector.Name()
	if interval == "monthly" {
		if dialect == "sqlite" {
			return "strftime('%Y-%m', date)"
		}
		return "DATE_FORMAT(date, '%Y-%m')"
	}
	if dialect == "sqlite" {
		return "strftime('%Y-%m-%d', date)"
	}
	return "DATE(date)"
}

// GetRevenueTimeline retrieves the revenue timeline within a specified date range and interval.
// It aggregates the total revenue from both prescription sales and regular sales by the given interval.
// The interval can be "daily" or "monthly", determining the granularity of the timeline.
// Results are returned as a slice of timelineQueryResult with each element representing the date and the total revenue for that date.
// The function returns an error if the query execution fails.
func (r *repository) GetRevenueTimeline(start, end time.Time, interval string) ([]timelineQueryResult, error) {
	var results []timelineQueryResult
	dateCol := "transaction_date"
	query := r.db.Raw(`
        SELECT T.date, SUM(T.value) as value
        FROM (
            SELECT DATE(`+dateCol+`) as date, SUM(total_amount) as value
            FROM prescription_sales
            WHERE `+dateCol+` BETWEEN ? AND ? AND deleted_at IS NULL
            GROUP BY DATE(`+dateCol+`)
            UNION ALL
            SELECT DATE(`+dateCol+`) as date, SUM(total_pay) as value
            FROM sales_regulars
            WHERE `+dateCol+` BETWEEN ? AND ? AND deleted_at IS NULL
            GROUP BY DATE(`+dateCol+`)
        ) AS T
        GROUP BY T.date
        ORDER BY T.date ASC`, start, end, start, end)

	err := query.Scan(&results).Error
	return results, err
}

// GetExpenseTimeline retrieves the expense timeline within a specified date range and interval.
// It aggregates the total expenses from the Expense table by the given interval.
// The interval can be "daily" or "monthly", determining the granularity of the timeline.
// Results are returned as a slice of timelineQueryResult with each element representing the date and the total expense for that date.
// The function returns an error if the query execution fails.
func (r *repository) GetExpenseTimeline(start, end time.Time, interval string) ([]timelineQueryResult, error) {
	var results []timelineQueryResult
	dateFormat := getDateFormatString(r.db, interval)

	err := r.db.Model(&expense.Expense{}).
		Select(fmt.Sprintf("%s as date, SUM(amount) as value", dateFormat)).
		Where("date BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

// GetTopSellingProducts retrieves the top selling products within a specified date range and limit.
// The result is a slice of TopProduct, sorted in descending order of total revenue.
// The function returns an error if the query execution fails.
func (r *repository) GetTopSellingProducts(start, end time.Time, limit int) ([]TopProduct, error) {
	var results []TopProduct

	query := r.db.Raw(`
        SELECT
            product_id,
            product_name,
            product_code,
            SUM(quantity) as total_quantity,
            SUM(revenue) as total_revenue
        FROM (
            SELECT
                s.product_id, 
                pi.item_name as product_name,
                pi.item_code as product_code,
                pi.quantity,
                pi.sub_total as revenue
            FROM prescription_items pi
            JOIN prescription_sales ps ON pi.prescription_sale_id = ps.id
            JOIN stocks s ON pi.stock_id = s.id
            WHERE ps.transaction_date BETWEEN ? AND ? AND pi.deleted_at IS NULL
            UNION ALL
            SELECT
                sri.product_id, 
                sri.product_name,    
                sri.product_code, 
                sri.qty as quantity,
                sri.sub_total as revenue  
            FROM sales_regular_items sri
            JOIN sales_regulars sr ON sri.sales_regular_id = sr.id
            WHERE sr.transaction_date BETWEEN ? AND ? AND sri.deleted_at IS NULL
        ) as combined_sales
        GROUP BY product_id, product_name, product_code
        ORDER BY total_revenue DESC
        LIMIT ?`, start, end, start, end, limit)

	err := query.Scan(&results).Error
	return results, err
}
