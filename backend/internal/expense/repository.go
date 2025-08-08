package expense

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(expense *Expense) (*Expense, error)
	FindAll() ([]Expense, error)
	FindByID(id uint) (*Expense, error)
	Update(expense *Expense) (*Expense, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(expense *Expense) (*Expense, error) {
	err := r.db.Create(expense).Error
	if err != nil {
		return nil, err
	}
	r.db.Preload("ExpenseType").First(expense, expense.ID)
	return expense, nil
}

func (r *repository) FindAll() ([]Expense, error) {
	var expenses []Expense
	err := r.db.Preload("ExpenseType").Order("date DESC").Find(&expenses).Error
	return expenses, err
}

func (r *repository) FindByID(id uint) (*Expense, error) {
	var expense Expense
	if err := r.db.Preload("ExpenseType").First(&expense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return &expense, nil
}

func (r *repository) Update(expense *Expense) (*Expense, error) {
	err := r.db.Save(expense).Error
	if err != nil {
		return nil, err
	}
	r.db.Preload("ExpenseType").First(expense, expense.ID)
	return expense, nil
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&Expense{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}
