package repository

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
)

type TransaksiRepository interface {
	Create(transaksi *model.Transaksi) error
	FindAll() ([]model.Transaksi, error)
	DeleteTransaksi(id uint) error
}

type transaksiRepository struct{}

func NewTransaksiRepository() TransaksiRepository {
	return &transaksiRepository{}
}

func (r *transaksiRepository) Create(transaksi *model.Transaksi) error {
	return config.DB.Create(transaksi).Error
}

func (r *transaksiRepository) FindAll() ([]model.Transaksi, error) {
	var transaksis []model.Transaksi
	err := config.DB.Find(&transaksis).Error
	return transaksis, err
}
func (r *transaksiRepository) DeleteTransaksi(id uint) error {
	var transaksi model.Transaksi
	if err := config.DB.First(&transaksi, id).Error; err != nil {
		return err
	}
	return config.DB.Delete(&transaksi).Error
}
