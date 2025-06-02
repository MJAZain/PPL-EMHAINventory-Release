package stock

import (
	"errors"
	"go-gin-auth/config"

	"gorm.io/gorm"
)

type Repository interface {
	GetProductStockById(id uint) (*Stock, error)
	UpdateProductStock(id uint, quantity int, isAdd bool) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository() *repository {
	return &repository{db: config.DB}
}

func (r *repository) GetProductStockById(id uint) (*Stock, error) {
	productStock := &Stock{}
	if err := r.db.Where("product_id = ?", id).First(productStock).Error; err != nil {
		return nil, err
	}
	return productStock, nil
}

func (r *repository) UpdateProductStock(id uint, quantity int, isAdd bool) error {
	productStock := &Stock{}
	if err := r.db.Where("product_id = ?", id).First(productStock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			productStock = &Stock{
				ProductID: id,
				Quantity:  quantity,
			}
			if err := r.db.Create(productStock).Error; err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if isAdd {
		productStock.Quantity += quantity
	} else {
		productStock.Quantity -= quantity
	}

	if err := r.db.Save(productStock).Error; err != nil {
		return err
	}

	return nil
}
