package outgoingProducts

import (
	"errors"
	"go-gin-auth/config"

	"gorm.io/gorm"
)

type Repository interface {
	Create(outgoingProduct *OutgoingProduct, details []OutgoingProductDetail) error
	GetAll() ([]OutgoingProduct, error)
	GetByID(id uint) (*OutgoingProduct, error)
	GetDetailsByOutgoingProductID(outgoingProductID uint) ([]OutgoingProductDetail, error)
	GetDetailByOutgoingProductID(outgoingProductID uint) (*OutgoingProductDetail, error)
	Update(id uint, outgoingProduct *OutgoingProduct) error
	UpdateDetails(details []OutgoingProductDetail) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository() *repository {
	return &repository{db: config.DB}
}

func (r *repository) Create(outgoingProduct *OutgoingProduct, details []OutgoingProductDetail) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(outgoingProduct).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range details {
		details[i].OutgoingProductID = outgoingProduct.ID
		if err := tx.Create(&details[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *repository) GetAll() ([]OutgoingProduct, error) {
	var outgoingProducts []OutgoingProduct
	if err := r.db.Find(&outgoingProducts).Error; err != nil {
		return nil, err
	}

	// Mengambil total amount untuk setiap outgoing product
	for i := range outgoingProducts {
		var details []OutgoingProductDetail
		if err := r.db.Where("outgoing_product_id = ?", outgoingProducts[i].ID).Find(&details).Error; err != nil {
			return nil, err
		}

		var totalAmount float64
		for _, detail := range details {
			totalAmount += detail.Total
		}
		outgoingProducts[i].TotalAmount = totalAmount
	}

	return outgoingProducts, nil
}

func (r *repository) GetByID(id uint) (*OutgoingProduct, error) {
	var outgoingProduct OutgoingProduct
	if err := r.db.First(&outgoingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outgoing product not found")
		}
		return nil, err
	}

	// Mengambil total amount
	var details []OutgoingProductDetail
	if err := r.db.Where("outgoing_product_id = ?", id).Find(&details).Error; err != nil {
		return nil, err
	}

	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.Total
	}
	outgoingProduct.TotalAmount = totalAmount

	return &outgoingProduct, nil
}

func (r *repository) GetDetailsByOutgoingProductID(outgoingProductID uint) ([]OutgoingProductDetail, error) {
	var details []OutgoingProductDetail
	if err := r.db.Where("outgoing_product_id = ?", outgoingProductID).Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *repository) Update(id uint, outgoingProduct *OutgoingProduct) error {
	var existingProduct OutgoingProduct
	if err := r.db.First(&existingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("outgoing product not found")
		}
		return err
	}

	// Update fields
	existingProduct.Date = outgoingProduct.Date
	existingProduct.Customer = outgoingProduct.Customer
	existingProduct.NoFaktur = outgoingProduct.NoFaktur
	existingProduct.PaymentStatus = outgoingProduct.PaymentStatus

	return r.db.Save(&existingProduct).Error
}

func (r *repository) UpdateDetails(details []OutgoingProductDetail) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, detail := range details {
		var existingDetail OutgoingProductDetail
		if err := tx.Where("outgoing_product_id = ? AND product_id = ?", detail.OutgoingProductID, detail.ProductID).First(&existingDetail).Error; err != nil {
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
	if err := tx.Where("outgoing_product_id = ?", id).Delete(&OutgoingProductDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Kemudian hapus outgoing product
	if err := tx.Delete(&OutgoingProduct{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *repository) GetDetailByOutgoingProductID(outgoingProductID uint) (*OutgoingProductDetail, error) {
	var detail OutgoingProductDetail
	if err := r.db.Where("outgoing_product_id = ?", outgoingProductID).First(&detail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outgoing product detail not found")
		}
		return nil, err
	}
	return &detail, nil
}
