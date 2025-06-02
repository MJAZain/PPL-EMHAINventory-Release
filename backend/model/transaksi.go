package model

import (
	"fmt"
	"math/rand"
	"time"
)

// Transaksi merepresentasikan data pembelian obat
type Transaksi struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	NomorTransaksi   string    `gorm:"uniqueIndex;not null" json:"nomor_transaksi"`
	ObatID           uint      `gorm:"not null" json:"obat_id"`
	JumlahObat       int       `gorm:"not null" json:"jumlah_obat"`
	TanggalPembelian time.Time `gorm:"not null" json:"tanggal_pembelian"`
	TotalHarga       float64   `gorm:"not null" json:"total_harga"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName menentukan nama tabel di database
func (Transaksi) TableName() string {
	return "transaksi"
}

// GenerateNomorTransaksi membuat nomor invoice unik, contoh: INV-20250506-XYZ123
func GenerateNomorTransaksi() string {
	timestamp := time.Now().Format("20060102")
	suffix := randomString(6)
	return fmt.Sprintf("INV-%s-%s", timestamp, suffix)
}

// randomString menghasilkan string acak untuk suffix invoice
func randomString(length int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
