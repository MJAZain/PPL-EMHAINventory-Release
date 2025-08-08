package sales

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalesRegularHandler struct {
	service SalesRegularService
}

func NewSalesRegularHandler(service SalesRegularService) *SalesRegularHandler {
	return &SalesRegularHandler{service: service}
}

// ✅ GET /api/sales/regular?limit=10&offset=0
func (h *SalesRegularHandler) GetAll(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	data, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": total,
	})
}

// ✅ GET /api/sales/regular/:id
func (h *SalesRegularHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	data, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// ✅ POST /api/sales/regular
func (h *SalesRegularHandler) Create(c *gin.Context) {
	var req SalesRegularRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat transaksi", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// ✅ PUT /api/sales/regular/:id
func (h *SalesRegularHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var req SalesRegularRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload tidak valid", "error": err.Error()})
		return
	}

	data, err := h.service.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate transaksi", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// ✅ DELETE /api/sales/regular/:id (optional)
func (h *SalesRegularHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus transaksi", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil dihapus"})
}
