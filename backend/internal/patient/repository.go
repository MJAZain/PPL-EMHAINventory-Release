package patient

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(patient *Patient) (*Patient, error)
	GetAll(searchQuery string) ([]Patient, error)
	GetByID(id uint) (*Patient, error)
	Update(id uint, patient *Patient) (*Patient, error)
	Delete(id uint) error
	FindActiveByIdentityNumber(identityNumber string) (*Patient, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindActiveByIdentityNumber(identityNumber string) (*Patient, error) {
	if identityNumber == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var patient Patient
	err := r.db.Where("LOWER(identity_number) = LOWER(?) AND status = ?", identityNumber, "Aktif").First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *repository) Create(patient *Patient) (*Patient, error) {
	if err := r.db.Create(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

func (r *repository) GetAll(searchQuery string) ([]Patient, error) {
	var patients []Patient
	query := r.db.Model(&Patient{})

	if searchQuery != "" {
		searchPattern := "%" + strings.ToLower(searchQuery) + "%"
		query = query.Where("LOWER(full_name) LIKE ? OR LOWER(identity_number) LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Order("id DESC").Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *repository) GetByID(id uint) (*Patient, error) {
	var patient Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &patient, nil
}

func (r *repository) Update(id uint, patientData *Patient) (*Patient, error) {
	if err := r.db.Model(&Patient{ID: id}).Updates(patientData).Error; err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

func (r *repository) Delete(id uint) error {
	result := r.db.Model(&Patient{}).Where("id = ?", id).Update("status", "Nonaktif")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
