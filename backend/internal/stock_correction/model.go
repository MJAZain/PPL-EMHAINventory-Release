package stock_correction

import "time"

type StockCorrection struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProductID         uint      `gorm:"not null;index" json:"product_id" form:"product_id"`
	ProductName       string    `gorm:"-" json:"product_name,omitempty"`
	OldStock          int       `gorm:"not null" json:"old_stock"`
	NewStock          int       `gorm:"not null" json:"new_stock" form:"new_stock"`
	Difference        int       `gorm:"not null" json:"difference"`
	Reason            string    `gorm:"type:varchar(255);not null" json:"reason" form:"reason"`
	CorrectionDate    time.Time `gorm:"not null" json:"correction_date"`
	CorrectionOfficer string    `gorm:"type:varchar(255);not null" json:"correction_officer"`
	Notes             string    `gorm:"type:text" json:"notes,omitempty" form:"notes"`
}
