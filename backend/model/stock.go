package model

import (
	"time"
)

type ProductStock struct {
	StockID        int       `json:"stock_id" gorm:"primaryKey;autoIncrement"`
	ProductID      string    `json:"product_id" gorm:"index"`
	CurrentStock   int       `json:"current_stock" gorm:"default:0"`
	LastOpnameDate time.Time `json:"last_opname_date"`
	LastUpdated    time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	// Product        product.Product `json:"product" gorm:"foreignKey:ProductID"`
}

// // ProductStockInfo adalah DTO untuk menampilkan informasi produk dan stok
// type ProductStockInfo struct {
// 	ProductID      string    `json:"product_id"`
// 	Name           string    `json:"name"`
// 	Category       string    `json:"category"`
// 	Unit           string    `json:"unit"`
// 	SystemStock    int       `json:"system_stock"`
// 	MinStock       int       `json:"min_stock"`
// 	LastOpnameDate time.Time `json:"last_opname_date"`
// }
