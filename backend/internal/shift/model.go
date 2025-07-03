package shift

import "time"

type Shift struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	ShiftDate        time.Time  `gorm:"type:date;not null" json:"shift_date"`
	OpeningOfficerID uint       `gorm:"not null" json:"opening_officer_id,omitempty" form:"opening_officer_id"`
	OpeningOfficer   string     `gorm:"-" json:"opening_officer_name,omitempty"`
	OpeningTime      time.Time  `gorm:"not null" json:"opening_time"`
	OpeningBalance   float64    `gorm:"type:decimal(14,2);not null" json:"opening_balance,omitempty" form:"opening_balance"`
	ClosingOfficerID uint       `json:"closing_officer_id,omitempty" form:"closing_officer_id"`
	ClosingOfficer   *string    `gorm:"-" json:"closing_officer_name,omitempty"`
	ClosingTime      *time.Time `json:"closing_time,omitempty"`
	ClosingBalance   *float64   `gorm:"type:decimal(14,2)" json:"closing_balance,omitempty" form:"closing_balance"`
	TotalSales       *float64   `gorm:"type:decimal(14,2)" json:"total_sales,omitempty" form:"total_sales"`
	ManualCorrection *float64   `gorm:"type:decimal(14,2)" json:"manual_correction,omitempty" form:"manual_correction"`
	Notes            string     `gorm:"type:text" json:"notes,omitempty" form:"notes"`
	Status           string     `gorm:"type:varchar(20);not null" json:"status"`
}
