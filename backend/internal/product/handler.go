package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{service: *NewProductService()}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to get products", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   products,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var productDTO Product
	err := c.Bind(&productDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to bind request",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, _ := c.Get("user_id")
	currentUserID := uint(currentUserIDFloat.(float64))
	productDTO.CreatedBy = currentUserID

	product, err := h.service.CreateProduct(productDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product created successfully",
		"data":    product,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to convert id",
			"error":   err.Error(),
		})
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Product not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product found",
		"data":    product,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to convert id",
			"error":   err.Error(),
		})
		return
	}

	var productDTO Product
	err = c.Bind(&productDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to bind request",
			"error":   err.Error(),
		})
		return
	}

	product, err := h.service.UpdateProduct(uint(id), productDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product updated successfully",
		"data":    product,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to convert id",
			"error":   err.Error(),
		})
		return
	}

	err = h.service.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product deleted successfully",
	})
}
