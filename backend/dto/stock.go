// DTO untuk menyesuaikan stok produk
package dto

import "time"

type StockAdjustmentRequest struct {
	ActualStock    int       `json:"actual_stock" binding:"required"`
	AdjustmentNote string    `json:"adjustment_note"`
	OpnameDate     time.Time `json:"opname_date" binding:"required"`
	PerformedBy    string    `json:"performed_by" binding:"required"`
}

// DTO untuk histori penyesuaian stok
type StockAdjustmentHistory struct {
	AdjustmentID          string    `json:"adjustment_id"`
	ProductID             string    `json:"product_id"`
	Name                  string    `json:"name"`
	PreviousStock         int       `json:"previous_stock"`
	ActualStock           int       `json:"actual_stock"`
	Discrepancy           int       `json:"discrepancy"`
	DiscrepancyPercentage float64   `json:"discrepancy_percentage"`
	AdjustmentNote        string    `json:"adjustment_note"`
	OpnameDate            time.Time `json:"opname_date"`
	PerformedBy           string    `json:"performed_by"`
}

// DTO untuk selisih stok yang signifikan
type StockDiscrepancy struct {
	ProductID             string    `json:"product_id"`
	Name                  string    `json:"name"`
	Category              string    `json:"category"`
	PreviousStock         int       `json:"previous_stock"`
	ActualStock           int       `json:"actual_stock"`
	Discrepancy           int       `json:"discrepancy"`
	DiscrepancyPercentage float64   `json:"discrepancy_percentage"`
	Flag                  string    `json:"flag"`
	OpnameDate            time.Time `json:"opname_date"`
	PerformedBy           string    `json:"performed_by"`
}

type ProductStockResponse struct {
	Name            string            `json:"name"`
	Code            string            `json:"code"`
	StockBuffer     int               `json:"stock_buffer"`
	StorageLocation string            `json:"storage_location"`
	Category        CategorySimpleDTO `json:"category"`
	Unit            UnitSimpleDTO     `json:"unit"`
}

type CategorySimpleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UnitSimpleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductSimple struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Code       string `json:"code"`
	// Barcode    string `json:"barcode"`
	// CategoryID uint   `json:"category_id"`
}

type StockOpnameDetailResponse struct {
	ID                     uint          `json:"detail_id"`
	QtySystem              int           `json:"system_stock"`
	QtyReal                int           `json:"actual_stock"`
	Discrepancy            int           `json:"discrepancy"`
	Discrepancy_percentage int           `json:"discrepancy_percentage"`
	Adjustment_note        string        `json:"adjustment_note"`
	Performed_by           string        `json:"performed_by"`
	Performed_at           time.Time     `json:"performed_at"`
	Product                ProductSimple `json:"product"`
}

type StockOpnameResponse struct {
	OpnameId        string                      `json:"opname_id"`
	OpnameDate      time.Time                   `json:"opname_date"`
	StartTime       time.Time                   `json:"start_time"`
	EndTime         time.Time                   `json:"end_time"`
	Status          string                      `json:"status"`
	Notes           string                      `json:"notes"`
	JenisStokOpname string                      `json:"jenis_stok_opname"`
	FlagActive      bool                        `json:"FlagActive"`
	CreatedBy       string                      `json:"created_by"`
	Details         []StockOpnameDetailResponse `json:"details"`
}
