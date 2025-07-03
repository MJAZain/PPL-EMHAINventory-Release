package patient

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

func (h *Handler) CreatePatient(c *gin.Context) {
	var input Patient
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	newPatient, err := h.service.CreatePatient(&input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrIdentityExists) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal membuat pasien", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusCreated, "Pasien berhasil dibuat", nil, newPatient)
}

func (h *Handler) GetAllPatients(c *gin.Context) {
	searchQuery := c.Query("search")
	patients, err := h.service.GetAllPatients(searchQuery)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil daftar pasien", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Daftar pasien berhasil diambil", nil, patients)
}

func (h *Handler) GetPatientByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}
	patient, err := h.service.GetPatientByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data pasien", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Detail pasien berhasil diambil", nil, patient)
}

func (h *Handler) UpdatePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}
	var input Patient
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}
	updatedPatient, err := h.service.UpdatePatient(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else if errors.Is(err, ErrIdentityExists) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal memperbarui data pasien", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Data pasien berhasil diperbarui", nil, updatedPatient)
}

func (h *Handler) DeletePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter ID tidak valid", err.Error(), nil)
		return
	}
	if err := h.service.DeletePatient(uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menonaktifkan pasien", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Pasien berhasil dinonaktifkan", nil, nil)
}
