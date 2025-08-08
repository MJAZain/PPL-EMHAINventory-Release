package expense

import (
	"errors"
	"go-gin-auth/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateExpense(c *gin.Context) {
	var input Expense
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}

	newExpense, err := h.service.CreateExpense(&input)
	if err != nil {
		if errors.Is(err, ErrInvalidAmount) || errors.Is(err, ErrInvalidTypeID) {
			utils.Respond(c, http.StatusBadRequest, "Failed to create expense", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "An internal error occurred", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusCreated, "Expense created successfully", nil, newExpense)
}

func (h *Handler) GetAllExpenses(c *gin.Context) {
	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve expenses", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Expenses retrieved successfully", nil, expenses)
}

func (h *Handler) GetExpenseByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	expense, err := h.service.GetExpenseByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrExpenseNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve expense", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense detail retrieved successfully", nil, expense)
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input Expense
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}

	updatedExpense, err := h.service.UpdateExpense(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrExpenseNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusBadRequest, "Failed to update expense", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense updated successfully", nil, updatedExpense)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteExpense(uint(id)); err != nil {
		if errors.Is(err, ErrExpenseNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Failed to delete expense", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense deleted successfully", nil, nil)
}
