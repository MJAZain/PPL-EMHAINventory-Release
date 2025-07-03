package doctor

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(doctor *Doctor) (*Doctor, error)
	GetAll(searchQuery string) ([]Doctor, error)
	GetByID(id uint) (*Doctor, error)
	Update(id uint, doctor *Doctor) (*Doctor, error)
	Delete(id uint) error
	GetAllActive() ([]Doctor, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllActive() ([]Doctor, error) {
	var doctors []Doctor
	if err := r.db.Where("status = ?", "Aktif").Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

func (r *repository) Create(doctor *Doctor) (*Doctor, error) {
	if err := r.db.Create(doctor).Error; err != nil {
		return nil, err
	}
	return doctor, nil
}

func (r *repository) GetAll(searchQuery string) ([]Doctor, error) {
	var doctors []Doctor
	query := r.db.Model(&Doctor{})
	if err := query.Order("id DESC").Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

func (r *repository) GetByID(id uint) (*Doctor, error) {
	var doctor Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &doctor, nil
}

func (r *repository) Update(id uint, doctorData *Doctor) (*Doctor, error) {
	if err := r.db.Model(&Doctor{ID: id}).Updates(doctorData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id uint) error {
	result := r.db.Model(&Doctor{}).Where("id = ?", id).Update("status", "Nonaktif")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
