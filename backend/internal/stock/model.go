package stock

import (
	"time"
)

type Stock struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ProductID    uint       `gorm:"not null;comment:ID Produk" json:"product_id"`
	ExpiryDate   *time.Time `gorm:"comment:Tanggal kedaluwarsa" json:"expiry_date"`            // untuk expiry tracking
	Quantity     int        `gorm:"not null;comment:Kuantitas" json:"quantity"`                // stok tersisa
	MinimumStock int        `gorm:"default:0;comment:Batas minimum stok" json:"minimum_stock"` // untuk deteksi stok kritis
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
