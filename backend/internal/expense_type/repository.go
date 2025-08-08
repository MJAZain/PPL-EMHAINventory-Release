package expense_type

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(expenseType *ExpenseType) (*ExpenseType, error)
	FindAll() ([]ExpenseType, error)
	FindByID(id uint) (*ExpenseType, error)
	Update(expenseType *ExpenseType) (*ExpenseType, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(expenseType *ExpenseType) (*ExpenseType, error) {
	err := r.db.Create(expenseType).Error
	return expenseType, err
}

func (r *repository) FindAll() ([]ExpenseType, error) {
	var expenseTypes []ExpenseType
	err := r.db.Order("name ASC").Find(&expenseTypes).Error
	return expenseTypes, err
}

func (r *repository) FindByID(id uint) (*ExpenseType, error) {
	var expenseType ExpenseType
	if err := r.db.First(&expenseType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseTypeNotFound
		}
		return nil, err
	}
	return &expenseType, nil
}

func (r *repository) Update(expenseType *ExpenseType) (*ExpenseType, error) {
	err := r.db.Save(expenseType).Error
	return expenseType, err
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&ExpenseType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExpenseTypeNotFound
	}
	return nil
}
