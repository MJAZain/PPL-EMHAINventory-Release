package product

import (
	"errors"
	"go-gin-auth/config"
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	"go-gin-auth/internal/drug_category"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/unit"
	"log"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: config.DB}
}

func (r *ProductRepository) GetProductByID(id uint) (Product, error) {
	var product Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return product, errors.New("product not found")
	}

	var category category.Category
	if err := r.db.First(&category, product.CategoryID).Error; err == nil {
		product.Category = category
	}
	var unit unit.Unit
	if err := r.db.First(&unit, product.UnitID).Error; err == nil {
		product.Unit = unit
	}

	var storageLocation storagelocation.StorageLocation
	if err := r.db.First(&storageLocation, product.StorageLocationID).Error; err == nil {
		product.StorageLocation = storageLocation
	}

	var brand brand.Brand
	if err := r.db.First(&brand, product.BrandID).Error; err == nil {
		product.Brand = brand
	}

	var drugCategory drug_category.DrugCategory
	if err := r.db.First(&drugCategory, product.DrugCategoryID).Error; err == nil {
		product.DrugCategory = drugCategory
	}

	return product, nil
}

func (r *ProductRepository) isCodeUsed(code string, id uint) bool {
	var count int64
	r.db.Model(&Product{}).
		Where("code = ?", code).
		Where("id != ?", id).
		Where("deleted_at IS NULL").
		Count(&count)
	return count > 0
}

func (r *ProductRepository) GetProducts() ([]Product, error) {
	var products []Product
	err := r.db.Find(&products).Error
	if err != nil {
		return products, errors.New("failed to retrieve products")
	}
	for i := range products {
		var category category.Category
		if err := r.db.First(&category, products[i].CategoryID).Error; err == nil {
			products[i].Category = category
		}
		var unit unit.Unit
		if err := r.db.First(&unit, products[i].UnitID).Error; err == nil {
			products[i].Unit = unit
		}
		var storageLocation storagelocation.StorageLocation
		log.Println("Retrieving storage location for product:", products[i].StorageLocationID)
		if err := r.db.First(&storageLocation, products[i].StorageLocationID).Error; err == nil {
			products[i].StorageLocation = storageLocation
		}
		var brand brand.Brand
		if err := r.db.First(&brand, products[i].BrandID).Error; err == nil {
			products[i].Brand = brand
		}

		var drugCategory drug_category.DrugCategory
		if err := r.db.First(&drugCategory, products[i].DrugCategoryID).Error; err == nil {
			products[i].DrugCategory = drugCategory
		}
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(product Product) (Product, error) {
	if r.isCodeUsed(product.Code, 0) {
		return product, errors.New("product code is already used")
	}

	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	var category category.Category
	if err := r.db.First(&category, product.CategoryID).Error; err == nil {
		product.Category = category
	}
	var unit unit.Unit
	if err := r.db.First(&unit, product.UnitID).Error; err == nil {
		product.Unit = unit
	}
	return product, nil
}

func (r *ProductRepository) UpdateProduct(id uint, product Product) (Product, error) {
	var existingProduct Product
	err := r.db.First(&existingProduct, id).Error
	if err != nil {
		return product, errors.New("product not found")
	}
	if r.isCodeUsed(product.Code, id) {
		return product, errors.New("product code is already used")
	}
	err = r.db.Model(&existingProduct).Updates(product).Error
	if err != nil {
		return product, errors.New("failed to update product")
	}

	var category category.Category
	if err := r.db.First(&category, existingProduct.CategoryID).Error; err == nil {
		existingProduct.Category = category
	}
	var unit unit.Unit
	if err := r.db.First(&unit, existingProduct.UnitID).Error; err == nil {
		existingProduct.Unit = unit
	}
	return existingProduct, nil
}

func (r *ProductRepository) DeleteProduct(id uint) error {
	var product Product
	err := r.db.Delete(&product, id).Error
	if err != nil {
		return errors.New("failed to delete product")
	}
	return nil
}
