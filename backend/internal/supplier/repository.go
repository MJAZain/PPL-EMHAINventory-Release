package supplier

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(supplier *Supplier) (*Supplier, error)
	GetAll(searchQuery string) ([]Supplier, error)
	GetByID(id uint) (*Supplier, error)
	Update(id uint, supplier *Supplier) (*Supplier, error)
	Delete(id uint) error
	FindActiveByName(name string) (*Supplier, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(supplier *Supplier) (*Supplier, error) {
	if err := r.db.Create(supplier).Error; err != nil {
		return nil, err
	}
	return supplier, nil
}

func (r *repository) GetAll(searchQuery string) ([]Supplier, error) {
	var suppliers []Supplier
	query := r.db.Model(&Supplier{})

	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		query = query.Where("name LIKE ?", searchPattern)
	}

	if err := query.Order("id DESC").Find(&suppliers).Error; err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *repository) GetByID(id uint) (*Supplier, error) {
	var supplier Supplier
	if err := r.db.First(&supplier, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &supplier, nil
}

func (r *repository) Update(id uint, supplierData *Supplier) (*Supplier, error) {
	if err := r.db.Model(&Supplier{ID: id}).Updates(supplierData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id uint) error {
	result := r.db.Model(&Supplier{}).Where("id = ?", id).Update("status", "Nonaktif")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) FindActiveByName(name string) (*Supplier, error) {
	var supplier Supplier
	err := r.db.Where("name = ? AND status = ?", name, "Aktif").First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}
