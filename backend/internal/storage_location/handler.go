package storagelocation

import (
	"go-gin-auth/pkg/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StorageLocationHandler struct {
	StorageLocationService *StorageLocationService
}

func NewStorageLocationHandler(service *StorageLocationService) *StorageLocationHandler {
	return &StorageLocationHandler{
		StorageLocationService: service,
	}
}

func (h *StorageLocationHandler) CreateStorageLocation(c *gin.Context) {
	var input StorageLocation

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

	sl, err := h.StorageLocationService.CreateStorageLocation(input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "kode lokasi penyimpanan sudah ada" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"status":  statusCode,
			"message": "Gagal membuat lokasi penyimpanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Lokasi Penyimpanan berhasil dibuat",
		"data":    sl,
	})
}

func (h *StorageLocationHandler) GetStorageLocations(c *gin.Context) {
	page, limit, _ := pagination.GetPaginationParams(c)
	search := c.DefaultQuery("search", "")

	storageLocations, totalData, err := h.StorageLocationService.GetStorageLocations(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal mendapatkan daftar lokasi penyimpanan",
			"error":   err.Error(),
		})
		return
	}

	paginatedResult := pagination.CreatePaginationResult(storageLocations, totalData, page, limit)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Daftar Lokasi Penyimpanan berhasil diambil",
		"data":    paginatedResult,
	})
}

func (h *StorageLocationHandler) GetStorageLocationByID(c *gin.Context) {
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

	slID := uint(id)
	sl, err := h.StorageLocationService.GetStorageLocationByID(slID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Lokasi Penyimpanan tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Lokasi Penyimpanan berhasil diambil",
		"data":    sl,
	})
}

func (h *StorageLocationHandler) UpdateStorageLocation(c *gin.Context) {
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
	slID := uint(id)

	var input StorageLocation
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

	updatedSL, err := h.StorageLocationService.UpdateStorageLocation(slID, input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "kode lokasi penyimpanan sudah ada" {
			statusCode = http.StatusConflict
		} else if err.Error() == "lokasi penyimpanan tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"status":  statusCode,
			"message": "Gagal memperbarui lokasi penyimpanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Lokasi Penyimpanan berhasil diperbarui",
		"data":    updatedSL,
	})
}

func (h *StorageLocationHandler) DeleteStorageLocation(c *gin.Context) {
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
	slID := uint(id)

	currentUserIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Pengguna belum terautentikasi",
		})
		return
	}
	currentUserID := uint(currentUserIDFloat.(float64))

	var slToDelete StorageLocation
	slToDelete.DeletedBy = currentUserID
	slToDelete.UpdatedBy = currentUserID

	if err := h.StorageLocationService.DeleteStorageLocation(slID, slToDelete); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Gagal menghapus lokasi penyimpanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Lokasi Penyimpanan berhasil dihapus",
	})
}
