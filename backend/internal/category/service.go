package category

import (
	"errors"
)

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repository: NewCategoryRepository()}
}

func (s *CategoryService) CreateCategory(category Category) (Category, error) {
	if category.Name == "" {
		return Category{}, errors.New("name is required")
	}
	return s.repository.CreateCategory(category)
}

func (s *CategoryService) GetCategories() ([]Category, error) {
	return s.repository.GetCategories()
}

func (s *CategoryService) GetCategoryByID(ID uint) (Category, error) {
	return s.repository.GetCategoryByID(ID)
}

func (s *CategoryService) UpdateCategory(ID uint, category Category) (Category, error) {
	if category.Name == "" {
		return Category{}, errors.New("name is required")
	}
	return s.repository.UpdateCategory(ID, category)
}

func (s *CategoryService) DeleteCategory(ID uint, category Category) error {
	return s.repository.DeleteCategory(ID, category)
}
