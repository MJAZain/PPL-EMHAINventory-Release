package supplier

import (
	"errors"
	"go-gin-auth/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSupplier(c *gin.Context) {
	var input Supplier
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	newSupplier, err := h.service.CreateSupplier(&input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal membuat supplier", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusCreated, "Supplier berhasil dibuat", nil, newSupplier)
}

func (h *Handler) GetAllSuppliers(c *gin.Context) {
	searchQuery := c.Query("search")

	suppliers, err := h.service.GetAllSuppliers(searchQuery)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil daftar supplier", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Daftar supplier berhasil diambil", nil, suppliers)
}

func (h *Handler) GetSupplierByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}

	supplier, err := h.service.GetSupplierByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "Supplier tidak ditemukan", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, "Detail supplier berhasil diambil", nil, supplier)
}

func (h *Handler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", nil, nil)
		return
	}

	var input Supplier
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	updatedSupplier, err := h.service.UpdateSupplier(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Respond(c, http.StatusNotFound, "Supplier tidak ditemukan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal memperbarui supplier", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Supplier berhasil diperbarui", nil, updatedSupplier)
}

func (h *Handler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", nil, nil)
		return
	}

	if err := h.service.DeleteSupplier(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Respond(c, http.StatusNotFound, "Supplier tidak ditemukan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menonaktifkan supplier", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Supplier berhasil dinonaktifkan", nil, nil)
}
