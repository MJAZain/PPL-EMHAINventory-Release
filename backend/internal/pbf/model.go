// Models
package pbf

import (
	"go-gin-auth/internal/product"
	"go-gin-auth/internal/supplier"
	"go-gin-auth/model"
	"time"
)

type IncomingPBF struct {
	ID              uint                `json:"id" gorm:"primaryKey"`
	OrderNumber     string              `json:"order_number" gorm:"not null"`
	OrderDate       time.Time           `json:"order_date" gorm:"not null"`
	ReceiptDate     time.Time           `json:"receipt_date" gorm:"not null"`
	TransactionCode string              `json:"transaction_code" gorm:"unique;not null"`
	SupplierID      uint                `json:"supplier_id" gorm:"not null"`
	InvoiceNumber   string              `json:"invoice_number" gorm:"not null"`
	TransactionType string              `json:"transaction_type" gorm:"type:varchar(20);default:'Cash'"`
	PaymentDueDate  *time.Time          `json:"payment_due_date"`
	UserID          uint                `json:"user_id" gorm:"not null"`
	AdditionalNotes string              `json:"additional_notes"`
	TotalPurchase   float64             `json:"total_purchase" gorm:"not null"`
	PaymentStatus   string              `json:"payment_status" gorm:"type:varchar(20);default:'Belum Lunas'"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Details         []IncomingPBFDetail `json:"details" gorm:"foreignKey:IncomingPBFID"`
	Supplier        supplier.Supplier   `json:"supplier" gorm:"foreignKey:SupplierID"`
	User            model.User          `json:"user" gorm:"foreignKey:UserID"`
}

type IncomingPBFDetail struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	IncomingPBFID uint            `json:"incoming_pbf_id" gorm:"not null"`
	ProductID     uint            `json:"product_id" gorm:"not null"`
	ProductCode   string          `json:"product_code" gorm:"not null"`
	ProductName   string          `json:"product_name" gorm:"not null"`
	Unit          string          `json:"unit" gorm:"not null"`
	Quantity      int             `json:"quantity" gorm:"not null"`
	PurchasePrice float64         `json:"purchase_price" gorm:"not null"`
	TotalPrice    float64         `json:"total_price" gorm:"not null"`
	BatchNumber   string          `json:"batch_number"`
	ExpiryDate    *time.Time      `json:"expiry_date"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Product       product.Product `json:"product" gorm:"foreignKey:ProductID"`
}
