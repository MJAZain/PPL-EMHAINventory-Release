// Models
package pbf

// Request/Response DTOs
type CreateIncomingPBFRequest struct {
	OrderNumber     string                           `json:"order_number" validate:"required"`
	OrderDate       string                           `json:"order_date" validate:"required"`
	ReceiptDate     string                           `json:"receipt_date" validate:"required"`
	SupplierID      uint                             `json:"supplier_id" validate:"required"`
	InvoiceNumber   string                           `json:"invoice_number" validate:"required"`
	TransactionType string                           `json:"transaction_type" validate:"required,oneof=Cash Kredit Konsinyasi"`
	PaymentDueDate  *string                          `json:"payment_due_date"`
	UserID          uint                             `json:"user_id" validate:"required"`
	AdditionalNotes string                           `json:"additional_notes"`
	PaymentStatus   string                           `json:"payment_status" validate:"oneof=Lunas 'Belum Lunas'"`
	Details         []CreateIncomingPBFDetailRequest `json:"details" validate:"required,min=1"`
}

type CreateIncomingPBFDetailRequest struct {
	ProductID     uint    `json:"product_id" validate:"required"`
	Quantity      int     `json:"quantity" validate:"required,min=1"`
	PurchasePrice float64 `json:"purchase_price" validate:"required,min=0"`
	BatchNumber   string  `json:"batch_number"`
	ExpiryDate    *string `json:"expiry_date"`
}
