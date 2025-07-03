package stock_correction

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(correction *StockCorrection) (*StockCorrection, error)
	GetAll() ([]StockCorrection, error)
	GetByID(id uint) (*StockCorrection, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(correction *StockCorrection) (*StockCorrection, error) {
	if err := r.db.Create(correction).Error; err != nil {
		return nil, err
	}
	return correction, nil
}

func (r *repository) GetAll() ([]StockCorrection, error) {
	var corrections []StockCorrection
	if err := r.db.Order("id DESC").Find(&corrections).Error; err != nil {
		return nil, err
	}
	return corrections, nil
}

func (r *repository) GetByID(id uint) (*StockCorrection, error) {
	var correction StockCorrection
	if err := r.db.First(&correction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &correction, nil
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&StockCorrection{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
