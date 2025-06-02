package brand

import (
	"go-gin-auth/pkg/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BrandHandler struct {
	BrandService *BrandService
}

func NewBrandHandler(service *BrandService) *BrandHandler {
	return &BrandHandler{
		BrandService: service,
	}
}

func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var input Brand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Pengguna belum terautentikasi",
		})
		return
	}
	currentUserID := uint(currentUserIDFloat.(float64))
	input.CreatedBy = currentUserID
	input.UpdatedBy = currentUserID

	brand, err := h.BrandService.CreateBrand(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal membuat brand",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Brand berhasil dibuat",
		"data":    brand,
	})
}

func (h *BrandHandler) GetBrands(c *gin.Context) {
	page, limit, _ := pagination.GetPaginationParams(c)
	search := c.DefaultQuery("search", "")

	brands, totalData, err := h.BrandService.GetBrands(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mendapatkan daftar brand",
			"error":   err.Error(),
		})
		return
	}

	paginatedResult := pagination.CreatePaginationResult(brands, totalData, page, limit)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Daftar brand berhasil diambil",
		"data":    paginatedResult,
	})
}

func (h *BrandHandler) GetBrandByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}

	brandID := uint(id)
	brand, err := h.BrandService.GetBrandByID(brandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Brand tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Brand berhasil diambil",
		"data":    brand,
	})
}

func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}
	brandID := uint(id)

	var input Brand
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
		return
	}

	currentUserIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Pengguna belum terautentikasi",
		})
		return
	}
	currentUserID := uint(currentUserIDFloat.(float64))
	input.UpdatedBy = currentUserID

	updatedBrand, err := h.BrandService.UpdateBrand(brandID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal memperbarui brand",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Brand berhasil diperbarui",
		"data":    updatedBrand,
	})
}

func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter ID tidak valid",
			"error":   err.Error(),
		})
		return
	}
	brandID := uint(id)

	currentUserIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Pengguna belum terautentikasi",
		})
		return
	}
	currentUserID := uint(currentUserIDFloat.(float64))

	var brandToDelete Brand
	brandToDelete.DeletedBy = currentUserID
	brandToDelete.UpdatedBy = currentUserID

	if err := h.BrandService.DeleteBrand(brandID, brandToDelete); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghapus brand",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Brand berhasil dihapus",
	})
}
