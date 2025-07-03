package doctor

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

func (h *Handler) CreateDoctor(c *gin.Context) {
	var input Doctor
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	newDoctor, err := h.service.CreateDoctor(&input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrSTRExists) {
			utils.Respond(c, http.StatusBadRequest, "Input tidak valid atau nomor STR sudah digunakan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal membuat dokter", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusCreated, "Dokter berhasil dibuat", nil, newDoctor)
}

func (h *Handler) GetAllDoctors(c *gin.Context) {
	searchQuery := c.Query("search")

	doctors, err := h.service.GetAllDoctors(searchQuery)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Respond(c, http.StatusOK, "Dokter tidak ditemukan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil daftar dokter", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Daftar dokter berhasil diambil", nil, doctors)
}

func (h *Handler) GetDoctorByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}

	doctor, err := h.service.GetDoctorByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, "Dokter tidak ditemukan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil detail dokter", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Detail dokter berhasil diambil", nil, doctor)
}

func (h *Handler) UpdateDoctor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}

	var input Doctor
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	updatedDoctor, err := h.service.UpdateDoctor(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, "Dokter tidak ditemukan", err.Error(), nil)
		} else if errors.Is(err, ErrSTRExists) {
			utils.Respond(c, http.StatusBadRequest, "Nomor STR sudah digunakan oleh dokter lain yang aktif", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal memperbarui data dokter", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Dokter berhasil diperbarui", nil, updatedDoctor)
}

func (h *Handler) DeleteDoctor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}

	if err := h.service.DeleteDoctor(uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, "Dokter tidak ditemukan", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menonaktifkan dokter", err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "Dokter berhasil dinonaktifkan", nil, nil)
}
