package shift

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNotFound         = errors.New("shift tidak ditemukan")
	ErrInvalidInput     = errors.New("input tidak valid atau tidak lengkap")
	ErrShiftAlreadyOpen = errors.New("masih ada shift yang aktif, harap tutup shift sebelumnya")
	ErrShiftNotOpen     = errors.New("shift ini tidak dalam status Buka")
)

type Service interface {
	OpenShift(shift *Shift) (*Shift, error)
	CloseShift(id uint, closingData *Shift) (*Shift, error)
	UpdateShift(id uint, updateData *Shift) (*Shift, error)
	GetAllShifts() ([]Shift, error)
	GetShiftByID(id uint) (*Shift, error)
	DeleteShift(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) OpenShift(shift *Shift) (*Shift, error) {
	existing, err := s.repository.FindOpenShift()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrShiftAlreadyOpen
	}

	shift.ShiftDate = time.Now()
	shift.OpeningTime = time.Now()
	shift.Status = "Buka"

	return s.repository.Create(shift)
}

func (s *service) CloseShift(id uint, closingData *Shift) (*Shift, error) {
	shift, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if shift.Status != "Buka" {
		return nil, ErrShiftNotOpen
	}

	now := time.Now()
	shift.ClosingOfficer = closingData.ClosingOfficer
	shift.ClosingTime = &now
	shift.ClosingBalance = closingData.ClosingBalance
	shift.TotalSales = closingData.TotalSales
	shift.Status = "Tutup"
	shift.Notes = closingData.Notes

	return s.repository.Update(id, shift)
}

func (s *service) UpdateShift(id uint, updateData *Shift) (*Shift, error) {
	if _, err := s.repository.GetByID(id); err != nil {
		return nil, err
	}
	return s.repository.Update(id, updateData)
}

func (s *service) GetAllShifts() ([]Shift, error) {
	return s.repository.GetAll()
}

func (s *service) GetShiftByID(id uint) (*Shift, error) {
	return s.repository.GetByID(id)
}

func (s *service) DeleteShift(id uint) error {
	return s.repository.Delete(id)
}
