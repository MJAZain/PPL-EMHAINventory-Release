package stock

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type StockWithProduct struct {
	StockID      uint       `json:"stock_id"`
	ProductID    uint       `json:"product_id"`
	ProductName  string     `json:"product_name"`
	ProductCode  string     `json:"product_code"`
	Quantity     int        `json:"quantity"`
	MinStock     int        `json:"min_stock"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	SellingPrice float64    `json:"selling_price"`
}
type BatchStockDTO struct {
	ProductID   uint       `json:"product_id"`
	ProductName string     `json:"product_name"`
	BatchNumber string     `json:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	Quantity    int        `json:"quantity"` // total masuk per batch
	Source      string     `json:"source"`   // "PBF" / "NonPBF"
}
type LowStockProduct struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalStock  int    `json:"total_stock"`
	MinStock    int    `json:"min_stock"`
}
type ExpiringStock struct {
	ProductID   uint       `json:"product_id"`
	ProductName string     `json:"product_name"`
	BatchNumber string     `json:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	Quantity    int        `json:"quantity"`
}
type StockSummary struct {
	TotalProducts     int `json:"total_products"`
	OutOfStockCount   int `json:"out_of_stock_count"`
	LowStockCount     int `json:"low_stock_count"`
	ExpiringSoonCount int `json:"expiring_soon_count"`
}
type StockDetail struct {
	ProductID   uint       `json:"product_id"`
	ProductName string     `json:"product_name"`
	BatchNumber string     `json:"batch_number"`
	ExpiryDate  *time.Time `json:"expiry_date"`
	Quantity    int        `json:"quantity"`
	Source      string     `json:"source"` // "PBF" atau "NonPBF"
}

type StockService struct {
	DB *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{DB: db}
}

func (s *StockService) GetCurrentStocks() ([]StockWithProduct, error) {
	var result []StockWithProduct

	query := `
	SELECT 
		s.id AS stock_id,
		s.product_id,
		p.name AS product_name,
		p.code AS product_code,
		s.quantity,
		p.min_stock,
		s.expiry_date,
		p.selling_price
	FROM stocks s
	JOIN products p ON p.id = s.product_id
	ORDER BY p.name
	`

	err := s.DB.Raw(query).Scan(&result).Error
	return result, err
}

func (s *StockService) GetStockBatches(itemID *uint) ([]BatchStockDTO, error) {
	var results []BatchStockDTO
	args := []interface{}{}

	where := ""
	if itemID != nil {
		where = "WHERE d.product_id = ?"
		args = append(args, *itemID, *itemID) // << isi dua kali untuk kedua SELECT
	}

	query := fmt.Sprintf(`
		SELECT 
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.quantity) AS quantity,
			'PBF' AS source
		FROM incoming_pbf_details d
		JOIN products p ON p.id = d.product_id
		%s
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		UNION

		SELECT 
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.incoming_quantity) AS quantity,
			'NonPBF' AS source
		FROM incoming_non_pbf_details d
		JOIN products p ON p.id = d.product_id
		%s
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date
	`, where, where)

	err := s.DB.Raw(query, args...).Scan(&results).Error
	return results, err
}

func (s *StockService) GetLowStock() ([]Stock, error) {
	var stocks []Stock
	err := s.DB.
		Model(&Stock{}).
		Where("quantity < minimum_stock").
		Find(&stocks).Error
	return stocks, err
}

func (s *StockService) GetExpiringSoonStocks(months int) ([]ExpiringStock, error) {
	var results []ExpiringStock

	// Jika months kurang dari 1, set ke default 3
	if months < 1 {
		months = 3
	}

	// Query SQL dengan interval dinamis pakai $1 parameter
	query := `
		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.quantity) AS quantity
		FROM incoming_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.expiry_date IS NOT NULL
		  AND d.expiry_date <= NOW() + INTERVAL '? months'
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		UNION ALL

		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.incoming_quantity) AS quantity
		FROM incoming_non_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.expiry_date IS NOT NULL
		  AND d.expiry_date <= NOW() + INTERVAL '? months'
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		ORDER BY expiry_date ASC
	`

	// Karena PostgreSQL gak bisa parameter interval langsung,
	// kita format query string dengan fmt.Sprintf
	query = fmt.Sprintf(`
		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.quantity) AS quantity
		FROM incoming_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.expiry_date IS NOT NULL
		  AND d.expiry_date <= NOW() + INTERVAL '%d months'
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		UNION ALL

		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.incoming_quantity) AS quantity
		FROM incoming_non_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.expiry_date IS NOT NULL
		  AND d.expiry_date <= NOW() + INTERVAL '%d months'
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		ORDER BY expiry_date ASC
	`, months, months)

	err := s.DB.Raw(query).Scan(&results).Error
	return results, err
}

func (s *StockService) GetStockSummary() (StockSummary, error) {
	var summary StockSummary

	query := `
		SELECT
			COUNT(DISTINCT product_id) AS total_products,
			COUNT(CASE WHEN total_quantity = 0 THEN 1 END) AS out_of_stock_count,
			COUNT(CASE WHEN total_quantity < min_stock THEN 1 END) AS low_stock_count,
			COUNT(DISTINCT CASE WHEN has_expiring THEN product_id END) AS expiring_soon_count
		FROM (
			SELECT 
				s.product_id,
				SUM(s.quantity) AS total_quantity,
				p.min_stock,
				BOOL_OR(s.expiry_date IS NOT NULL AND s.expiry_date <= NOW() + INTERVAL '3 months') AS has_expiring
			FROM stocks s
			JOIN products p ON p.id = s.product_id
			GROUP BY s.product_id, p.min_stock
		) sub
	`

	err := s.DB.Raw(query).Scan(&summary).Error
	return summary, err
}

func (s *StockService) GetStockDetail(productID uint) ([]StockDetail, error) {
	var results []StockDetail

	query := `
		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.quantity) AS quantity,
			'PBF' AS source
		FROM incoming_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.product_id = ?
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		UNION ALL

		SELECT
			d.product_id,
			p.name AS product_name,
			d.batch_number,
			d.expiry_date,
			SUM(d.incoming_quantity) AS quantity,
			'NonPBF' AS source
		FROM incoming_non_pbf_details d
		JOIN products p ON p.id = d.product_id
		WHERE d.product_id = ?
		GROUP BY d.product_id, p.name, d.batch_number, d.expiry_date

		ORDER BY expiry_date ASC
	`

	err := s.DB.Raw(query, productID, productID).Scan(&results).Error
	return results, err
}
