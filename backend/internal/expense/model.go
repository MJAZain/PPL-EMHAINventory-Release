package expense

import (
	"go-gin-auth/internal/expense_type"
	"time"
)

type Expense struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ExpenseTypeID uint      `gorm:"not null;index" json:"expense_type_id" form:"expense_type_id"`
	Amount        float64   `gorm:"not null" json:"amount" form:"amount"`
	Description   string    `gorm:"type:text" json:"description,omitempty" form:"description"`
	Date          time.Time `gorm:"not null" json:"date" form:"date"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	ExpenseType expense_type.ExpenseType `gorm:"foreignKey:ExpenseTypeID" json:"expense_type"`
}
