package stock_correction

import (
	"errors"
	"go-gin-auth/internal/stock"
	"strings"
	"time"
)

var (
	ErrNotFound          = errors.New("data koreksi stok tidak ditemukan")
	ErrInvalidInput      = errors.New("input tidak valid atau tidak lengkap")
	ErrStockUpdateFailed = errors.New("gagal memperbarui data stok utama")
)

type Service interface {
	CreateCorrection(correction *StockCorrection, officerName string) (*StockCorrection, error)
	GetAllCorrections() ([]StockCorrection, error)
	GetCorrectionByID(id uint) (*StockCorrection, error)
	DeleteCorrection(id uint) error
}

type service struct {
	repository      Repository
	stockRepository stock.Repository
}

func NewService(repo Repository, stockRepo stock.Repository) Service {
	return &service{
		repository:      repo,
		stockRepository: stockRepo,
	}
}

func (s *service) CreateCorrection(correction *StockCorrection, officerName string) (*StockCorrection, error) {
	correction.Reason = strings.TrimSpace(correction.Reason)
	if correction.ProductID == 0 || correction.Reason == "" {
		return nil, ErrInvalidInput
	}

	currentStock, err := s.stockRepository.GetProductStockById(correction.ProductID)
	if err != nil {
		currentStock = &stock.Stock{ProductID: correction.ProductID, Quantity: 0}
	}

	correction.OldStock = currentStock.Quantity
	correction.Difference = correction.NewStock - correction.OldStock
	correction.CorrectionDate = time.Now()
	correction.CorrectionOfficer = officerName

	txErr := s.stockRepository.UpdateProductStock(correction.ProductID, correction.Difference, true)
	if txErr != nil {
		return nil, ErrStockUpdateFailed
	}

	return s.repository.Create(correction)
}

func (s *service) GetAllCorrections() ([]StockCorrection, error) {
	return s.repository.GetAll()
}

func (s *service) GetCorrectionByID(id uint) (*StockCorrection, error) {
	return s.repository.GetByID(id)
}

func (s *service) DeleteCorrection(id uint) error {
	return s.repository.Delete(id)
}
