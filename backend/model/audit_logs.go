package model

import "time"

type AuditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TableName   string    `json:"table_name"`
	RecordID    string    `json:"record_id"`
	Action      string    `json:"action"`      // INSERT, UPDATE, DELETE
	Description string    `json:"description"` // New field for description of the action
	ChangedBy   string    `json:"changed_by"`
	ChangedAt   time.Time `json:"changed_at"`
	BeforeData  string    `json:"before_data"` // JSON string
	AfterData   string    `json:"after_data"`  // JSON string
}
