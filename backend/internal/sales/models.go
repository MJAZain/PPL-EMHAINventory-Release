package sales

import (
	"time"

	"gorm.io/gorm"
)

type SalesRegular struct {
	ID              uint               `gorm:"primaryKey" json:"id"`
	SalesCode       string             `gorm:"uniqueIndex;not null" json:"sales_code"`
	TransactionDate time.Time          `gorm:"not null" json:"transaction_date"`
	CashierName     string             `gorm:"not null" json:"cashier_name"`
	CustomerName    *string            `json:"customer_name,omitempty"`
	CustomerContact *string            `json:"customer_contact,omitempty"`
	Description     *string            `json:"description,omitempty"`
	SubTotal        int                `gorm:"not null" json:"sub_total"`
	TotalDiscount   *int               `json:"total_discount,omitempty"`
	TotalPay        int                `gorm:"not null" json:"total_pay"`
	PaymentMethod   string             `gorm:"not null" json:"payment_method"`
	ShiftID         *uint              `json:"shift_id,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	DeletedAt       gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
	Items           []SalesRegularItem `gorm:"foreignKey:SalesRegularID" json:"items"`
}

type SalesRegularItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SalesRegularID uint           `gorm:"not null" json:"sales_regular_id"`
	ProductID      uint           `gorm:"not null" json:"product_id"`
	ProductCode    string         `gorm:"not null" json:"product_code"`
	ProductName    string         `gorm:"not null" json:"product_name"`
	Qty            int            `gorm:"not null" json:"qty"`
	Unit           string         `gorm:"not null" json:"unit"`
	UnitPrice      int            `gorm:"not null" json:"unit_price"`
	SubTotal       int            `gorm:"not null" json:"sub_total"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
type SalesRegularItemRequest struct {
	ProductID   uint   `json:"product_id"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	Qty         int    `json:"qty"`
	Unit        string `json:"unit"`
	UnitPrice   int    `json:"unit_price"`
	SubTotal    int    `json:"sub_total"`
}

type SalesRegularRequest struct {
	TransactionDate time.Time                 `json:"transaction_date"`
	CashierName     string                    `json:"cashier_name"`
	CustomerName    *string                   `json:"customer_name"`
	CustomerContact *string                   `json:"customer_contact"`
	Description     *string                   `json:"description"`
	SubTotal        int                       `json:"sub_total"`
	TotalDiscount   *int                      `json:"total_discount"`
	TotalPay        int                       `json:"total_pay"`
	PaymentMethod   string                    `json:"payment_method"`
	ShiftID         *uint                     `json:"shift_id"`
	Items           []SalesRegularItemRequest `json:"items"`
}
