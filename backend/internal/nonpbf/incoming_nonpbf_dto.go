package nonpbf

import "time"

type CreateIncomingNonPBFRequest struct {
	OrderNumber     string                        `json:"order_number" binding:"required"`
	OrderDate       time.Time                     `json:"order_date" binding:"required"`
	IncomingDate    time.Time                     `json:"incoming_date" binding:"required"`
	SupplierName    string                        `json:"supplier_name" binding:"required"`
	InvoiceNumber   string                        `json:"invoice_number" binding:"required"`
	TransactionType string                        `json:"transaction_type" binding:"required,oneof=Cash Kredit Konsinyasi"`
	PaymentDueDate  *time.Time                    `json:"payment_due_date"`
	OfficerName     string                        `json:"officer_name" binding:"required"`
	AdditionalNotes string                        `json:"additional_notes"`
	PaymentStatus   string                        `json:"payment_status" binding:"omitempty,oneof=Lunas 'Belum Lunas'"`
	UserID          uint                          `json:"user_id" binding:"required"`
	Details         []CreateIncomingDetailRequest `json:"details" binding:"required,dive"`
}

type CreateIncomingDetailRequest struct {
	ProductCode      string     `json:"product_code" binding:"required"`
	ProductName      string     `json:"product_name" binding:"required"`
	Unit             string     `json:"unit" binding:"required"`
	IncomingQuantity int        `json:"incoming_quantity" binding:"required,min=1"`
	PurchasePrice    float64    `json:"purchase_price" binding:"required,min=0.01"`
	BatchNumber      string     `json:"batch_number"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	ProductID        *uint      `json:"product_id"`
}

type UpdateIncomingNonPBFRequest struct {
	OrderNumber     string                        `json:"order_number"`
	OrderDate       time.Time                     `json:"order_date"`
	IncomingDate    time.Time                     `json:"incoming_date"`
	SupplierName    string                        `json:"supplier_name"`
	InvoiceNumber   string                        `json:"invoice_number"`
	TransactionType string                        `json:"transaction_type" binding:"omitempty,oneof=Cash Kredit Konsinyasi"`
	PaymentDueDate  *time.Time                    `json:"payment_due_date"`
	OfficerName     string                        `json:"officer_name"`
	AdditionalNotes string                        `json:"additional_notes"`
	PaymentStatus   string                        `json:"payment_status" binding:"omitempty,oneof=Lunas 'Belum Lunas'"`
	Details         []CreateIncomingDetailRequest `json:"details" binding:"dive"`
}
