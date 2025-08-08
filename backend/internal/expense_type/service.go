package expense_type

import (
	"errors"
	"strings"
)

var (
	ErrExpenseTypeNotFound = errors.New("expense type not found")
	ErrInvalidInput        = errors.New("invalid input: name cannot be empty")
)

type Service interface {
	CreateExpenseType(input *ExpenseType) (*ExpenseType, error)
	GetAllExpenseTypes() ([]ExpenseType, error)
	GetExpenseTypeByID(id uint) (*ExpenseType, error)
	UpdateExpenseType(id uint, input *ExpenseType) (*ExpenseType, error)
	DeleteExpenseType(id uint) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{repository: repo}
}

func (s *service) CreateExpenseType(input *ExpenseType) (*ExpenseType, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, ErrInvalidInput
	}
	return s.repository.Create(input)
}

func (s *service) GetAllExpenseTypes() ([]ExpenseType, error) {
	return s.repository.FindAll()
}

func (s *service) GetExpenseTypeByID(id uint) (*ExpenseType, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateExpenseType(id uint, input *ExpenseType) (*ExpenseType, error) {
	existingType, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, ErrInvalidInput
	}

	existingType.Name = input.Name
	return s.repository.Update(existingType)
}

func (s *service) DeleteExpenseType(id uint) error {
	_, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(id)
}
