package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryService CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		CategoryService: *NewCategoryService(),
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input Category

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

	category, err := h.CategoryService.CreateCategory(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Category created successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.CategoryService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to get categories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   categories,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	categoryID := uint(id)
	category, err := h.CategoryService.GetCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Category not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   category,
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	var input Category

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

	categoryID := uint(id)
	category, err := h.CategoryService.UpdateCategory(categoryID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	categoryID := uint(id)
	category, err := h.CategoryService.GetCategoryByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Category not found",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, _ := c.Get("user_id")
	category.DeletedBy = uint(currentUserIDFloat.(float64))

	if err := h.CategoryService.DeleteCategory(categoryID, category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete category",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Category deleted successfully",
	})
}
