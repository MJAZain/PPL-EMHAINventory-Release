package drug_category

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(category *DrugCategory) (*DrugCategory, error)
	GetAll() ([]DrugCategory, error)
	GetByID(id uint) (*DrugCategory, error)
	Update(id uint, category *DrugCategory) (*DrugCategory, error)
	Delete(id uint) error
	FindActiveByName(name string) (*DrugCategory, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindActiveByName(name string) (*DrugCategory, error) {
	var category DrugCategory
	err := r.db.Where("name = ? AND status = ?", name, "Aktif").First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) Create(category *DrugCategory) (*DrugCategory, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *repository) GetAll() ([]DrugCategory, error) {
	var categories []DrugCategory
	if err := r.db.Order("id DESC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) GetByID(id uint) (*DrugCategory, error) {
	var category DrugCategory
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *repository) Update(id uint, categoryData *DrugCategory) (*DrugCategory, error) {
	if err := r.db.Model(&DrugCategory{ID: id}).Updates(categoryData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id uint) error {
	result := r.db.Model(&DrugCategory{}).Where("id = ?", id).Update("status", "Nonaktif")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
