package model

type StockDiscrepancyFlag struct {
	FlagID          int     `json:"flag_id" gorm:"primaryKey;autoIncrement"`
	FlagName        string  `json:"flag_name" gorm:"not null"`
	MinPercentage   float64 `json:"min_percentage" gorm:"not null"`
	MaxPercentage   float64 `json:"max_percentage" gorm:"not null"`
	FlagColor       string  `json:"flag_color"`
	RequireApproval bool    `json:"require_approval" gorm:"default:false"`
}
