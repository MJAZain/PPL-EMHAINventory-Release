package drug_category

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("golongan obat tidak ditemukan")
	ErrInvalidInput = errors.New("input tidak valid, nama tidak boleh kosong")
	ErrNameExists   = errors.New("nama golongan obat yang aktif sudah ada")
)

type Service interface {
	CreateCategory(category *DrugCategory) (*DrugCategory, error)
	GetAllCategories() ([]DrugCategory, error)
	GetCategoryByID(id uint) (*DrugCategory, error)
	UpdateCategory(id uint, category *DrugCategory) (*DrugCategory, error)
	DeleteCategory(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) CreateCategory(category *DrugCategory) (*DrugCategory, error) {
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repository.FindActiveByName(category.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrNameExists
	}

	if category.Status == "" {
		category.Status = "Aktif"
	}

	return s.repository.Create(category)
}

func (s *service) UpdateCategory(id uint, category *DrugCategory) (*DrugCategory, error) {
	if _, err := s.repository.GetByID(id); err != nil {
		return nil, err
	}

	if category.Name != "" {
		category.Name = strings.TrimSpace(category.Name)
		existing, err := s.repository.FindActiveByName(category.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrNameExists
		}
	}

	return s.repository.Update(id, category)
}

func (s *service) GetAllCategories() ([]DrugCategory, error) {
	return s.repository.GetAll()
}

func (s *service) GetCategoryByID(id uint) (*DrugCategory, error) {
	return s.repository.GetByID(id)
}

func (s *service) DeleteCategory(id uint) error {
	if _, err := s.repository.GetByID(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}
