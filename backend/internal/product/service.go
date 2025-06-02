package product

import (
	"errors"
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/unit"
)

type ProductService struct {
	productRepo         ProductRepository
	categoryRepo        category.CategoryRepository
	unitRepo            unit.UnitRepository
	brandRepo           brand.BrandRepository
	storageLocationRepo storagelocation.StorageLocationRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo:         *NewProductRepository(),
		categoryRepo:        category.NewCategoryRepository(),
		unitRepo:            unit.NewUnitRepository(),
		brandRepo:           brand.NewBrandRepository(),
		storageLocationRepo: storagelocation.NewStorageLocationRepository(),
	}
}

func (s *ProductService) GetProductByID(id uint) (Product, error) {
	return s.productRepo.GetProductByID(id)
}

func (s *ProductService) GetProducts() ([]Product, error) {
	return s.productRepo.GetProducts()
}

func (s *ProductService) CreateProduct(product Product) (Product, error) {
	if _, err := s.categoryRepo.GetCategoryByID(product.CategoryID); err != nil {
		return product, errors.New("category not found")
	}

	if _, err := s.unitRepo.GetUnitByID(product.UnitID); err != nil {
		return product, errors.New("unit not found")
	}

	if _, err := s.brandRepo.GetBrandByID(product.BrandID); err != nil {
		return product, errors.New("brand not found")
	}

	if _, err := s.storageLocationRepo.GetStorageLocationByID(product.StorageLocationID); err != nil {
		return product, errors.New("storage location not found")
	}

	if err := validateProductFields(product); err != nil {
		return product, err
	}

	return s.productRepo.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(id uint, product Product) (Product, error) {
	if _, err := s.categoryRepo.GetCategoryByID(product.CategoryID); err != nil {
		return product, errors.New("category not found")
	}

	if _, err := s.unitRepo.GetUnitByID(product.UnitID); err != nil {
		return product, errors.New("unit not found")
	}

	if _, err := s.brandRepo.GetBrandByID(product.BrandID); err != nil {
		return product, errors.New("brand not found")
	}

	if _, err := s.storageLocationRepo.GetStorageLocationByID(product.StorageLocationID); err != nil {
		return product, errors.New("storage location not found")
	}

	return s.productRepo.UpdateProduct(id, product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

func validateProductFields(product Product) error {
	if product.Name == "" {
		return errors.New("name is required")
	}
	if product.Code == "" {
		return errors.New("code is required")
	}
	if product.Barcode == "" {
		return errors.New("barcode is required")
	}
	if product.PackageContent == 0 {
		return errors.New("package content is required")
	}
	if product.PurchasePrice == 0 {
		return errors.New("purchase price is required")
	}
	if product.SellingPrice == 0 {
		return errors.New("selling price is required")
	}
	return nil
}
