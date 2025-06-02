package unit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	UnitService UnitService
}

func NewUnitHandler() *UnitHandler {
	return &UnitHandler{
		UnitService: *NewUnitService(),
	}
}

func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var input Unit

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, _ := c.Get("user_id")
	currentUserID := uint(currentUserIDFloat.(float64))
	input.CreatedBy = currentUserID

	unit, err := h.UnitService.CreateUnit(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create unit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Unit created successfully",
		"data":    unit,
	})
}

func (h *UnitHandler) GetUnits(c *gin.Context) {
	units, err := h.UnitService.GetUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to get units",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   units,
	})
}

func (h *UnitHandler) GetUnitByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	unitID := uint(id)
	unit, err := h.UnitService.GetUnitByID(unitID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Unit not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   unit,
	})
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	var input Unit

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, _ := c.Get("user_id")
	input.UpdatedBy = uint(currentUserIDFloat.(float64))

	unitID := uint(id)
	unit, err := h.UnitService.UpdateUnit(unitID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update unit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Unit updated successfully",
		"data":    unit,
	})
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	unitID := uint(id)
	unit, err := h.UnitService.GetUnitByID(unitID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Unit not found",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, _ := c.Get("user_id")
	unit.DeletedBy = uint(currentUserIDFloat.(float64))

	if err := h.UnitService.DeleteUnit(unitID, unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete unit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Unit deleted successfully",
	})
}
