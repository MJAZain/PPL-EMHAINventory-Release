package outgoingProducts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandlerOutgoingProducts() *Handler {
	return &Handler{service: NewService()}
}

func (h *Handler) CreateOutgoingProduct(c *gin.Context) {
	var request struct {
		OutgoingProduct OutgoingProduct         `json:"outgoing_product"`
		Details         []OutgoingProductDetail `json:"details"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateOutgoingProduct(&request.OutgoingProduct, request.Details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membuat produk keluar",
			"error":   err.Error(),
		})
		return
	}

	product, err := h.service.GetOutgoingProductByID(request.OutgoingProduct.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Produk keluar berhasil dibuat tetapi gagal mengambil data",
			"error":   err.Error(),
		})
		return
	}

	details, err := h.service.GetOutgoingProductDetails(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Produk keluar berhasil dibuat tetapi gagal mengambil detail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Produk keluar berhasil dibuat",
		"data": gin.H{
			"outgoing_product": product,
			"details":          details,
		},
	})
}

func (h *Handler) GetAllOutgoingProducts(c *gin.Context) {
	products, err := h.service.GetAllOutgoingProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mendapatkan daftar produk keluar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   products,
	})
}

func (h *Handler) GetOutgoingProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	product, err := h.service.GetOutgoingProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Produk keluar tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	details, err := h.service.GetOutgoingProductDetails(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mendapatkan detail produk keluar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"outgoing_product": product,
			"details":          details,
		},
	})
}

func (h *Handler) UpdateOutgoingProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	var requestData struct {
		OutgoingProduct OutgoingProduct         `json:"outgoing_product"`
		Details         []OutgoingProductDetail `json:"details"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// 1. Update produk keluar
	if err := h.service.UpdateOutgoingProduct(uint(id), &requestData.OutgoingProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mengupdate produk keluar",
			"error":   err.Error(),
		})
		return
	}

	// 2. Update detail produk keluar jika ada
	if len(requestData.Details) > 0 {
		// Pastikan semua detail memiliki outgoing_product_id yang sama
		for i := range requestData.Details {
			requestData.Details[i].OutgoingProductID = uint(id)
		}

		if err := h.service.UpdateOutgoingProductDetails(requestData.Details); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengupdate detail produk keluar",
				"error":   err.Error(),
			})
			return
		}
	}

	// Ambil data produk keluar yang sudah diupdate
	product, err := h.service.GetOutgoingProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Produk keluar berhasil diupdate tetapi gagal mengambil data",
			"error":   err.Error(),
		})
		return
	}

	// Ambil detail produk keluar yang sudah diupdate
	updatedDetails, err := h.service.GetOutgoingProductDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Produk keluar berhasil diupdate tetapi gagal mengambil data detail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Produk keluar dan detailnya berhasil diupdate",
		"data": gin.H{
			"outgoing_product": product,
			"details":          updatedDetails,
		},
	})
}

// DeleteOutgoingProduct godoc
// @Summary Menghapus produk keluar
// @Description Menghapus produk keluar berdasarkan ID
// @Tags OutgoingProducts
// @Produce json
// @Param id path int true "ID Produk Keluar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /outgoing-products/{id} [delete]
func (h *Handler) DeleteOutgoingProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Pastikan produk keluar ada
	_, err = h.service.GetOutgoingProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Produk keluar tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.DeleteOutgoingProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghapus produk keluar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Produk keluar berhasil dihapus",
	})
}
