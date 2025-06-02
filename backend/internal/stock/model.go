package stock

import "time"

type Stock struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"not null;comment:ID Produk" json:"product_id"`
	Quantity       int       `gorm:"not null;comment:Kuantitas" json:"quantity"`
	LastOpnameDate time.Time `json:"last_opname_date"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
