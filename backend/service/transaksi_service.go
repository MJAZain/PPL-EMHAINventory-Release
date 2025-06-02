package service

import (
	"go-gin-auth/model"
	"go-gin-auth/repository"
	"time"
)

type TransaksiService interface {
	CreateTransaksi(input model.Transaksi) (*model.Transaksi, error)
	GetAllTransaksi() ([]model.Transaksi, error)
	DeleteTransaksi(id uint) error
}

type transaksiService struct {
	repo repository.TransaksiRepository
}

func NewTransaksiService(r repository.TransaksiRepository) TransaksiService {
	return &transaksiService{repo: r}
}

func (s *transaksiService) CreateTransaksi(input model.Transaksi) (*model.Transaksi, error) {
	input.NomorTransaksi = model.GenerateNomorTransaksi()
	input.TanggalPembelian = time.Now()

	err := s.repo.Create(&input)
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *transaksiService) GetAllTransaksi() ([]model.Transaksi, error) {
	return s.repo.FindAll()
}
func (s *transaksiService) DeleteTransaksi(id uint) error {
	return s.repo.DeleteTransaksi(id)
}
