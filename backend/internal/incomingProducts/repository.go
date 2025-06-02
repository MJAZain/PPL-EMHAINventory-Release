package incomingProducts

import (
	"errors"
	"go-gin-auth/config"

	"gorm.io/gorm"
)

type Repository interface {
	Create(incomingProduct *IncomingProduct, details []IncomingProductDetail) error
	GetAll() ([]IncomingProduct, error)
	GetByID(id uint) (*IncomingProduct, error)
	GetDetailsByIncomingProductID(incomingProductID uint) ([]IncomingProductDetail, error)
	GetDetailByIncomingProductID(incomingProductID uint) (*IncomingProductDetail, error)
	Update(id uint, incomingProduct *IncomingProduct) error
	UpdateDetails(details []IncomingProductDetail) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository() *repository {
	return &repository{db: config.DB}
}

func (r *repository) Create(incomingProduct *IncomingProduct, details []IncomingProductDetail) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(incomingProduct).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range details {
		details[i].IncomingProductID = incomingProduct.ID
		if err := tx.Create(&details[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) GetAll() ([]IncomingProduct, error) {
	var incomingProducts []IncomingProduct
	if err := r.db.Find(&incomingProducts).Error; err != nil {
		return nil, err
	}

	// Mengambil total amount untuk setiap incoming product
	for i := range incomingProducts {
		var details []IncomingProductDetail
		if err := r.db.Where("incoming_product_id = ?", incomingProducts[i].ID).Find(&details).Error; err != nil {
			return nil, err
		}

		var totalAmount float64
		for _, detail := range details {
			totalAmount += detail.Total
		}
		incomingProducts[i].TotalAmount = totalAmount
	}

	return incomingProducts, nil
}

func (r *repository) GetByID(id uint) (*IncomingProduct, error) {
	var incomingProduct IncomingProduct
	if err := r.db.First(&incomingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("incoming product not found")
		}
		return nil, err
	}

	// Mengambil total amount
	var details []IncomingProductDetail
	if err := r.db.Where("incoming_product_id = ?", id).Find(&details).Error; err != nil {
		return nil, err
	}

	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.Total
	}
	incomingProduct.TotalAmount = totalAmount

	return &incomingProduct, nil
}

func (r *repository) GetDetailsByIncomingProductID(incomingProductID uint) ([]IncomingProductDetail, error) {
	var details []IncomingProductDetail
	if err := r.db.Where("incoming_product_id = ?", incomingProductID).Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *repository) Update(id uint, incomingProduct *IncomingProduct) error {
	var existingProduct IncomingProduct
	if err := r.db.First(&existingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("incoming product not found")
		}
		return err
	}

	// Update fields
	existingProduct.Date = incomingProduct.Date
	existingProduct.Supplier = incomingProduct.Supplier
	existingProduct.NoFaktur = incomingProduct.NoFaktur
	existingProduct.PaymentStatus = incomingProduct.PaymentStatus

	return r.db.Save(&existingProduct).Error
}

func (r *repository) UpdateDetails(details []IncomingProductDetail) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, detail := range details {
		var existingDetail IncomingProductDetail
		if err := tx.Where("incoming_product_id = ? AND product_id = ?", detail.IncomingProductID, detail.ProductID).First(&existingDetail).Error; err != nil {
			tx.Rollback()
			return err
		}

		existingDetail.ProductID = detail.ProductID
		existingDetail.Quantity = detail.Quantity
		existingDetail.Price = detail.Price
		existingDetail.Total = detail.Total

		if err := tx.Save(&existingDetail).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) Delete(id uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Hapus detail terlebih dahulu
	if err := tx.Where("incoming_product_id = ?", id).Delete(&IncomingProductDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Kemudian hapus incoming product
	if err := tx.Delete(&IncomingProduct{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) GetDetailByIncomingProductID(incomingProductID uint) (*IncomingProductDetail, error) {
	var detail IncomingProductDetail
	if err := r.db.Where("incoming_product_id = ?", incomingProductID).First(&detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("incoming product detail not found")
		}
		return nil, err
	}
	return &detail, nil
}
