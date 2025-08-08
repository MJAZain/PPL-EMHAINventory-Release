package prescription

import (
	"go-gin-auth/internal/doctor"
	"go-gin-auth/internal/patient"
	"go-gin-auth/internal/shift"
	"go-gin-auth/internal/stock"
	"time"

	"gorm.io/gorm"
)

// PrescriptionSale represents the main prescription sale record
type PrescriptionSale struct {
	ID               uint          `json:"id" gorm:"primaryKey"`
	TransactionCode  string        `json:"transaction_code" gorm:"unique;not null"`
	PrescriptionNo   string        `json:"prescription_no" gorm:"not null"`
	PrescriptionDate time.Time     `json:"prescription_date"`
	DoctorID         uint          `json:"doctor_id"`
	Doctor           doctor.Doctor `json:"doctor" gorm:"foreignKey:DoctorID"`
	Clinic           string        `json:"clinic"`
	Diagnosis        string        `json:"diagnosis"`

	// Patient Info
	PatientID uint            `json:"patient_id"`
	Patient   patient.Patient `json:"patient" gorm:"foreignKey:PatientID"`

	// Transaction Info
	TransactionDate time.Time   `json:"transaction_date"`
	PaymentMethod   string      `json:"payment_method"` // Tunai/Transfer/QRIS/Kredit/BPJS/Jaminan
	DiscountPercent float64     `json:"discount_percent"`
	DiscountAmount  float64     `json:"discount_amount"`
	TotalAmount     float64     `json:"total_amount"`
	ShiftID         uint        `json:"shift_id"`
	Shift           shift.Shift `json:"shift" gorm:"foreignKey:ShiftID"`

	// Items with proper cascade delete
	Items []PrescriptionItem `json:"items" gorm:"foreignKey:PrescriptionSaleID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// PrescriptionItem represents individual items in a prescription sale
type PrescriptionItem struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	PrescriptionSaleID uint           `json:"prescription_sale_id" gorm:"not null;index"`
	StockID            uint           `json:"stock_id" gorm:"not null"`
	Stock              stock.Stock    `json:"stock" gorm:"foreignKey:StockID"`
	ItemCode           string         `json:"item_code"`
	ItemName           string         `json:"item_name"`
	Quantity           int            `json:"quantity" gorm:"not null;check:quantity > 0"`
	Unit               string         `json:"unit"`
	Price              float64        `json:"price" gorm:"not null;check:price >= 0"`
	SubTotal           float64        `json:"sub_total" gorm:"not null;check:sub_total >= 0"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (PrescriptionItem) TableName() string {
	return "prescription_items"
}

// CreatePrescriptionSaleRequest represents the request payload for creating/updating prescription sales
type CreatePrescriptionSaleRequest struct {
	PrescriptionNo   string                          `json:"prescription_no" binding:"required"`
	PrescriptionDate time.Time                       `json:"prescription_date" binding:"required"`
	DoctorID         uint                            `json:"doctor_id" binding:"required"`
	Clinic           string                          `json:"clinic"`
	Diagnosis        string                          `json:"diagnosis"`
	PatientID        uint                            `json:"patient_id" binding:"required"`
	TransactionDate  time.Time                       `json:"transaction_date" binding:"required"`
	PaymentMethod    string                          `json:"payment_method" binding:"required"`
	DiscountPercent  float64                         `json:"discount_percent"`
	DiscountAmount   float64                         `json:"discount_amount"`
	ShiftID          uint                            `json:"shift_id" binding:"required"`
	Items            []CreatePrescriptionItemRequest `json:"items" binding:"required,min=1"`
}

// CreatePrescriptionItemRequest represents individual item in the request
type CreatePrescriptionItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Code      string  `json:"code" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Unit      string  `json:"unit" binding:"required"`
	Price     float64 `json:"price" binding:"required,min=0"`
}
