package unit

import (
	"errors"
)

type UnitService struct {
	repository UnitRepository
}

func NewUnitService() *UnitService {
	return &UnitService{repository: NewUnitRepository()}
}

func (s *UnitService) CreateUnit(unit Unit) (Unit, error) {
	if unit.Name == "" {
		return Unit{}, errors.New("name is required")
	}
	return s.repository.CreateUnit(unit)
}

func (s *UnitService) GetUnits() ([]Unit, error) {
	return s.repository.GetUnits()
}

func (s *UnitService) GetUnitByID(ID uint) (Unit, error) {
	return s.repository.GetUnitByID(ID)
}

func (s *UnitService) UpdateUnit(ID uint, unit Unit) (Unit, error) {
	if unit.Name == "" {
		return Unit{}, errors.New("name is required")
	}
	return s.repository.UpdateUnit(ID, unit)
}

func (s *UnitService) DeleteUnit(ID uint, unit Unit) error {
	return s.repository.DeleteUnit(ID, unit)
}
