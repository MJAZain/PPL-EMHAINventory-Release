package brand

import (
	"errors"
)

type BrandService struct {
	repository BrandRepository
}

func NewBrandService(repo BrandRepository) *BrandService {
	return &BrandService{repository: repo}
}

func (s *BrandService) CreateBrand(brand Brand) (Brand, error) {
	if brand.Name == "" {
		return Brand{}, errors.New("nama wajib diisi")
	}
	return s.repository.CreateBrand(brand)
}

func (s *BrandService) GetBrands(page, limit int, search string) ([]Brand, int64, error) {
	return s.repository.GetBrands(page, limit, search)
}

func (s *BrandService) GetBrandByID(ID uint) (Brand, error) {
	return s.repository.GetBrandByID(ID)
}

func (s *BrandService) UpdateBrand(ID uint, brand Brand) (Brand, error) {
	if brand.Name == "" {
		return Brand{}, errors.New("nama wajib diisi untuk pembaruan")
	}
	return s.repository.UpdateBrand(ID, brand)
}

func (s *BrandService) DeleteBrand(ID uint, brand Brand) error {
	return s.repository.DeleteBrand(ID, brand)
}
