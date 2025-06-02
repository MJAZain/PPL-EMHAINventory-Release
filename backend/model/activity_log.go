package model

import "time"

type ActivityLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	FullName     string    `gorm:"not null" json:"full_name"`
	ActivityType string    `gorm:"size:50;not null" json:"activity_type"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	ActivityTime time.Time `gorm:"autoCreateTime" json:"activity_time"`
	IPAddress    string    `gorm:"size:50" json:"ip_address"`
	UserAgent    string    `gorm:"type:text" json:"user_agent"`
}
