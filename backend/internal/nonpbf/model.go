package nonpbf

import (
	"go-gin-auth/internal/product"
	"go-gin-auth/model"
	"time"

	"gorm.io/gorm"
)

type IncomingNonPBF struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	OrderNumber     string         `json:"order_number" gorm:"size:100;not null"`
	OrderDate       time.Time      `json:"order_date" gorm:"not null"`
	IncomingDate    time.Time      `json:"incoming_date" gorm:"not null"`
	TransactionCode string         `json:"transaction_code" gorm:"size:100;unique;not null"`
	SupplierName    string         `json:"supplier_name" gorm:"size:255;not null"`
	InvoiceNumber   string         `json:"invoice_number" gorm:"size:100;not null"`
	TransactionType string         `json:"transaction_type" gorm:"size:20;not null"` // Cash/Kredit/Konsinyasi
	PaymentDueDate  *time.Time     `json:"payment_due_date"`
	OfficerName     string         `json:"officer_name" gorm:"size:255;not null"`
	AdditionalNotes string         `json:"additional_notes" gorm:"type:text"`
	TotalPurchase   float64        `json:"total_purchase" gorm:"type:decimal(15,2);not null"`
	PaymentStatus   string         `json:"payment_status" gorm:"size:20;default:'Belum Lunas'"` // Lunas/Belum Lunas
	UserID          uint           `json:"user_id" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relations
	User                  model.User             `json:"user" gorm:"foreignKey:UserID"`
	IncomingNonPBFDetails []IncomingNonPBFDetail `json:"details" gorm:"foreignKey:IncomingNonPBFID"`
}

type IncomingNonPBFDetail struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	IncomingNonPBFID uint           `json:"incoming_nonpbf_id" gorm:"not null"`
	ProductCode      string         `json:"product_code" gorm:"size:100;not null"`
	ProductName      string         `json:"product_name" gorm:"size:255;not null"`
	Unit             string         `json:"unit" gorm:"size:50;not null"`
	IncomingQuantity int            `json:"incoming_quantity" gorm:"not null"`
	PurchasePrice    float64        `json:"purchase_price" gorm:"type:decimal(15,2);not null"`
	TotalPurchase    float64        `json:"total_purchase" gorm:"type:decimal(15,2);not null"`
	BatchNumber      string         `json:"batch_number" gorm:"size:100"`
	ExpiryDate       *time.Time     `json:"expiry_date"`
	ProductID        *uint          `json:"product_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relations
	IncomingNonPBF IncomingNonPBF  `json:"-" gorm:"foreignKey:IncomingNonPBFID"`
	Product        product.Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
