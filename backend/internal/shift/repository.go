package shift

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(shift *Shift) (*Shift, error)
	GetAll() ([]Shift, error)
	GetByID(id uint) (*Shift, error)
	Update(id uint, shift *Shift) (*Shift, error)
	Delete(id uint) error
	FindOpenShift() (*Shift, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOpenShift() (*Shift, error) {
	var shift Shift
	err := r.db.Where("status = ?", "Buka").First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *repository) Create(shift *Shift) (*Shift, error) {
	if err := r.db.Create(shift).Error; err != nil {
		return nil, err
	}
	return shift, nil
}

func (r *repository) GetAll() ([]Shift, error) {
	var shifts []Shift
	if err := r.db.Order("id DESC").Find(&shifts).Error; err != nil {
		return nil, err
	}
	return shifts, nil
}

func (r *repository) GetByID(id uint) (*Shift, error) {
	var shift Shift
	if err := r.db.First(&shift, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &shift, nil
}

func (r *repository) Update(id uint, shiftData *Shift) (*Shift, error) {
	if err := r.db.Model(&Shift{ID: id}).Updates(shiftData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&Shift{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
