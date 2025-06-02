package opname

import (
	"go-gin-auth/internal/product"
	"time"
)

type StockOpnameDetail struct {
	DetailID  int    `json:"detail_id" gorm:"primaryKey;autoIncrement"`
	OpnameID  string `json:"opname_id" gorm:"index"`
	ProductID uint   `json:"product_id" gorm:"index"`

	SystemStock           int             `json:"system_stock" gorm:"not null"`
	ActualStock           int             `json:"actual_stock" gorm:"not null"`
	Discrepancy           int             `json:"discrepancy" gorm:"not null"`
	DiscrepancyPercentage float64         `json:"discrepancy_percentage" gorm:"not null"`
	AdjustmentNote        string          `json:"adjustment_note"`
	PerformedBy           string          `json:"performed_by" gorm:"not null"`
	PerformedAt           time.Time       `json:"performed_at" gorm:"autoCreateTime"`
	Product               product.Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`
}

func (d *StockOpnameDetail) CalculateDiscrepancy() {
	d.Discrepancy = d.ActualStock - d.SystemStock
	if d.SystemStock == 0 {
		if d.ActualStock == 0 {
			d.DiscrepancyPercentage = 0
		} else {
			d.DiscrepancyPercentage = 100
		}
	} else {
		d.DiscrepancyPercentage = float64(d.Discrepancy) * 100.0 / float64(d.SystemStock)
	}
}

type StockOpnameStatus string

const (
	Draft      StockOpnameStatus = "draft"
	InProgress StockOpnameStatus = "in_progress"
	Completed  StockOpnameStatus = "completed"
	Canceled   StockOpnameStatus = "canceled"
)

type JenisStokOpname string

const (
	Regular JenisStokOpname = "Regular"
	Harian  JenisStokOpname = "Harian"
)

type StockOpname struct {
	// ID        uint                `gorm:"primaryKey" json:"id"`
	// UserID    uint                `gorm:"not null" json:"user_id"`
	// CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	// Details   []StockOpnameDetail `gorm:"foreignKey:StockOpnameID;constraint:OnDelete:CASCADE;" json:"details,omitempty"`
	OpnameID   string              `json:"opname_id" gorm:"primaryKey"`
	OpnameDate time.Time           `json:"opname_date" gorm:"not null"`
	StartTime  time.Time           `json:"start_time"`
	EndTime    time.Time           `json:"end_time"`
	Status     StockOpnameStatus   `json:"status" gorm:"default:'draft'"`
	Notes      string              `json:"notes"`
	Jenis      JenisStokOpname     `gorm:"column:jenis_stok_opname;type:varchar(20);not null" json:"jenis_stok_opname"`
	FlagActive bool                `gorm:"column:flag_active;default:true"`
	CreatedBy  string              `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
	Details    []StockOpnameDetail `json:"details" gorm:"foreignKey:OpnameID"`
}
