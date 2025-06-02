package adjustment

import (
	"time"
)

type AdjustmentType string

const (
	Opname     AdjustmentType = "opname"
	Damage     AdjustmentType = "damage"
	Expired    AdjustmentType = "expired"
	Correction AdjustmentType = "correction"
	Other      AdjustmentType = "other"
)

type StockAdjustment struct {
	AdjustmentID       string         `json:"adjustment_id" gorm:"primaryKey"`
	ProductID          string         `json:"product_id" gorm:"index"`
	PreviousStock      int            `json:"previous_stock" gorm:"not null"`
	AdjustedStock      int            `json:"adjusted_stock" gorm:"not null"`
	AdjustmentQuantity int            `json:"adjustment_quantity" gorm:"-"`
	AdjustmentType     AdjustmentType `json:"adjustment_type" gorm:"not null"`
	ReferenceID        string         `json:"reference_id"`
	AdjustmentNote     string         `json:"adjustment_note"`
	AdjustmentDate     time.Time      `json:"adjustment_date" gorm:"autoCreateTime"`
	PerformedBy        string         `json:"performed_by" gorm:"not null"`
	//Product            product.Product `json:"product" gorm:"foreignKey:ProductID"`
}

func (a *StockAdjustment) CalculateAdjustmentQuantity() {
	a.AdjustmentQuantity = a.AdjustedStock - a.PreviousStock
}
