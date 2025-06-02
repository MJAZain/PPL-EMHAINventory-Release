package dto

type StockOpnameDetailRequest struct {
	ObatID    uint `json:"obat_id" binding:"required"`
	StokFisik int  `json:"stok_fisik" binding:"required"`
}

type StockOpnameRequest struct {
	UserID  uint                       `json:"user_id" binding:"required"`
	Details []StockOpnameDetailRequest `json:"details" binding:"required,dive"`
}

type AddProductRequest struct {
	ProductID string `json:"product_id" binding:"required"`
}

type RecordStockRequest struct {
	ActualStock int    `json:"actual_stock" binding:"required"`
	Note        string `json:"note"`
}
