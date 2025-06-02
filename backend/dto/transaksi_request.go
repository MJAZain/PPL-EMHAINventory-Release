package dto

type CreateTransaksiRequest struct {
	ObatID     uint    `json:"obat_id" binding:"required"`
	JumlahObat int     `json:"jumlah_obat" binding:"required"`
	TotalHarga float64 `json:"total_harga" binding:"required"`
	UserID     uint    `json:"user_id" binding:"required"`
}
