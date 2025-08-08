package expense

import (
	"errors"
	"go-gin-auth/internal/expense_type"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
	ErrInvalidAmount   = errors.New("invalid input: amount must be greater than zero")
	ErrInvalidTypeID   = errors.New("invalid input: provided expense type does not exist")
)

type Service interface {
	CreateExpense(input *Expense) (*Expense, error)
	GetAllExpenses() ([]Expense, error)
	GetExpenseByID(id uint) (*Expense, error)
	UpdateExpense(id uint, input *Expense) (*Expense, error)
	DeleteExpense(id uint) error
}

type service struct {
	repository         Repository
	expenseTypeService expense_type.Service
}

func NewService(repo Repository, expenseTypeService expense_type.Service) Service {
	return &service{
		repository:         repo,
		expenseTypeService: expenseTypeService,
	}
}

func (s *service) CreateExpense(input *Expense) (*Expense, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if _, err := s.expenseTypeService.GetExpenseTypeByID(input.ExpenseTypeID); err != nil {
		return nil, ErrInvalidTypeID
	}

	return s.repository.Create(input)
}

func (s *service) GetAllExpenses() ([]Expense, error) {
	return s.repository.FindAll()
}

func (s *service) GetExpenseByID(id uint) (*Expense, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateExpense(id uint, input *Expense) (*Expense, error) {
	existingExpense, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if input.ExpenseTypeID != existingExpense.ExpenseTypeID {
		if _, err := s.expenseTypeService.GetExpenseTypeByID(input.ExpenseTypeID); err != nil {
			return nil, ErrInvalidTypeID
		}
	}

	existingExpense.ExpenseTypeID = input.ExpenseTypeID
	existingExpense.Amount = input.Amount
	existingExpense.Description = input.Description
	existingExpense.Date = input.Date

	return s.repository.Update(existingExpense)
}

func (s *service) DeleteExpense(id uint) error {
	_, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(id)
}
