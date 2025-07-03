package drug_category

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

func (h *Handler) CreateCategory(c *gin.Context) {
	var input DrugCategory
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}
	newCategory, err := h.service.CreateCategory(&input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrNameExists) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal membuat golongan obat", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusCreated, "Golongan obat berhasil dibuat", nil, newCategory)
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil daftar golongan obat", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Daftar golongan obat berhasil diambil", nil, categories)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data golongan obat", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Detail golongan obat berhasil diambil", nil, category)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input DrugCategory
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}
	updatedCategory, err := h.service.UpdateCategory(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else if errors.Is(err, ErrNameExists) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal memperbarui data golongan obat", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Data golongan obat berhasil diperbarui", nil, updatedCategory)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menonaktifkan golongan obat", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Golongan obat berhasil dinonaktifkan", nil, nil)
}
