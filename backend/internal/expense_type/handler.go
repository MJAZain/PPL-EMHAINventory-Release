package expense_type

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

func (h *Handler) CreateExpenseType(c *gin.Context) {
	var input ExpenseType
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}

	newType, err := h.service.CreateExpenseType(&input)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Failed to create expense type", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusCreated, "Expense type created successfully", nil, newType)
}

func (h *Handler) GetAllExpenseTypes(c *gin.Context) {
	types, err := h.service.GetAllExpenseTypes()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve expense types", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Expense types retrieved successfully", nil, types)
}

func (h *Handler) GetExpenseTypeByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	expenseType, err := h.service.GetExpenseTypeByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrExpenseTypeNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Failed to retrieve expense type", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense type detail retrieved successfully", nil, expenseType)
}

func (h *Handler) UpdateExpenseType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input ExpenseType
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}

	updatedType, err := h.service.UpdateExpenseType(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrExpenseTypeNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusBadRequest, "Failed to update expense type", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense type updated successfully", nil, updatedType)
}

func (h *Handler) DeleteExpenseType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteExpenseType(uint(id)); err != nil {
		if errors.Is(err, ErrExpenseTypeNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Failed to delete expense type", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Expense type deleted successfully", nil, nil)
}
