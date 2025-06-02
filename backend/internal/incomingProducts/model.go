package incomingProducts

type IncomingProduct struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	Date          string  `gorm:"type:date;not null;comment:Tanggal" json:"date"`
	Supplier      string  `gorm:"type:varchar(100);not null;comment:Supplier" json:"supplier"`
	NoFaktur      string  `gorm:"type:varchar(100);not null;comment:No Faktur" json:"no_faktur"`
	PaymentStatus string  `gorm:"type:varchar(100);not null;comment:Status Pembayaran" json:"payment_status"`
	TotalAmount   float64 `gorm:"-" json:"total_amount"`
}

type IncomingProductDetail struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	IncomingProductID uint    `gorm:"not null;comment:ID Produk Masuk" json:"incoming_product_id"`
	ProductID         uint    `gorm:"not null;comment:ID Produk" json:"product_id"`
	Quantity          int     `gorm:"not null;comment:Kuantitas" json:"quantity"`
	Price             float64 `gorm:"not null;comment:Harga" json:"price"`
	Total             float64 `gorm:"not null;comment:Total" json:"total"`
}
