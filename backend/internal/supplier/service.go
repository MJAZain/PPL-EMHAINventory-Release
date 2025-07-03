package supplier

import (
	"errors"
	"go-gin-auth/internal/location"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("supplier tidak ditemukan")
	ErrInvalidInput    = errors.New("input tidak valid atau tidak lengkap")
	ErrInvalidLocation = errors.New("kombinasi ID provinsi dan kota/kabupaten tidak valid")
	ErrNameExists      = errors.New("nama supplier yang aktif sudah ada")
)

type Service interface {
	CreateSupplier(supplier *Supplier) (*Supplier, error)
	GetAllSuppliers(searchQuery string) ([]Supplier, error)
	GetSupplierByID(id uint) (*Supplier, error)
	UpdateSupplier(id uint, supplier *Supplier) (*Supplier, error)
	DeleteSupplier(id uint) error
}

type service struct {
	repository      Repository
	locationService location.Service
}

func NewService(repository Repository, locationService location.Service) Service {
	return &service{repository, locationService}
}

func (s *service) enrichSupplierData(supplier *Supplier) {
	if supplier != nil {
		supplier.Province, supplier.City = s.locationService.GetLocationNames(supplier.ProvinceID, supplier.CityID)
	}
}

func (s *service) CreateSupplier(supplier *Supplier) (*Supplier, error) {
	supplier.Name = strings.TrimSpace(supplier.Name)
	if supplier.Name == "" || supplier.ProvinceID == "" || supplier.CityID == "" {
		return nil, ErrInvalidInput
	}

	existing, err := s.repository.FindActiveByName(supplier.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrNameExists
	}

	valid, _ := s.locationService.ValidateLocation(supplier.ProvinceID, supplier.CityID)
	if !valid {
		return nil, ErrInvalidLocation
	}

	if supplier.Status == "" {
		supplier.Status = "Aktif"
	}

	newSupplier, err := s.repository.Create(supplier)
	if err != nil {
		return nil, err
	}

	s.enrichSupplierData(newSupplier)
	return newSupplier, nil
}

func (s *service) UpdateSupplier(id uint, supplier *Supplier) (*Supplier, error) {
	if _, err := s.repository.GetByID(id); err != nil {
		return nil, err
	}

	if supplier.Name != "" {
		supplier.Name = strings.TrimSpace(supplier.Name)
		existing, err := s.repository.FindActiveByName(supplier.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrNameExists
		}
	}

	if supplier.ProvinceID != "" && supplier.CityID != "" {
		valid, _ := s.locationService.ValidateLocation(supplier.ProvinceID, supplier.CityID)
		if !valid {
			return nil, ErrInvalidLocation
		}
	}

	updatedSupplier, err := s.repository.Update(id, supplier)
	if err != nil {
		return nil, err
	}

	s.enrichSupplierData(updatedSupplier)
	return updatedSupplier, nil
}

func (s *service) GetSupplierByID(id uint) (*Supplier, error) {
	supplier, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	s.enrichSupplierData(supplier)
	return supplier, nil
}

func (s *service) GetAllSuppliers(searchQuery string) ([]Supplier, error) {
	suppliers, err := s.repository.GetAll(searchQuery)
	if err != nil {
		return nil, err
	}
	for i := range suppliers {
		s.enrichSupplierData(&suppliers[i])
	}
	return suppliers, nil
}

func (s *service) DeleteSupplier(id uint) error {
	if _, err := s.repository.GetByID(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}
