// Handler
package prescription

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrescriptionSaleHandler struct {
	service *PrescriptionSaleService
}

func NewPrescriptionSaleHandler(service *PrescriptionSaleService) *PrescriptionSaleHandler {
	return &PrescriptionSaleHandler{service: service}
}

func (h *PrescriptionSaleHandler) GetAll(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	sales, total, err := h.service.GetAll(pageInt, limitInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":  sales,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

func (h *PrescriptionSaleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	sale, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Prescription sale not found"})
		return
	}

	c.JSON(200, gin.H{"data": sale})
}

func (h *PrescriptionSaleHandler) Create(c *gin.Context) {
	var req CreatePrescriptionSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.service.Create(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": sale})
}

func (h *PrescriptionSaleHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req CreatePrescriptionSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.service.Update(uint(id), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": sale})
}

func (h *PrescriptionSaleHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := h.service.Delete(uint(id))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Prescription sale deleted successfully"})
}
