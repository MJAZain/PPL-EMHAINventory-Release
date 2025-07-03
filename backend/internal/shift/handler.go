package shift

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

func (h *Handler) OpenShift(c *gin.Context) {
	var input Shift
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	userID, _ := c.Get("user_id")
	officerName, _ := c.Get("full_name")
	input.OpeningOfficerID = uint(userID.(float64))
	input.OpeningOfficer = officerName.(string)

	newShift, err := h.service.OpenShift(&input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrShiftAlreadyOpen) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal membuka shift", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusCreated, "Shift berhasil dibuka", nil, newShift)
}

func (h *Handler) CloseShift(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input Shift
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	userID, _ := c.Get("user_id")
	officerName, _ := c.Get("full_name")
	input.ClosingOfficerID = uint(userID.(float64))
	if nameStr, ok := officerName.(string); ok {
		input.ClosingOfficer = &nameStr
	} else {
		input.ClosingOfficer = nil
	}

	closedShift, err := h.service.CloseShift(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrShiftNotOpen) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menutup shift", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Shift berhasil ditutup", nil, closedShift)
}

func (h *Handler) UpdateShift(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input Shift
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}
	updatedShift, err := h.service.UpdateShift(uint(id), &input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal memperbarui shift", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Shift berhasil diperbarui", nil, updatedShift)
}

func (h *Handler) GetAllShifts(c *gin.Context) {
	shifts, err := h.service.GetAllShifts()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil daftar shift", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Daftar shift berhasil diambil", nil, shifts)
}

func (h *Handler) GetShiftByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	shift, err := h.service.GetShiftByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data shift", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Detail shift berhasil diambil", nil, shift)
}

func (h *Handler) DeleteShift(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteShift(uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menghapus shift", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Shift berhasil dihapus", nil, nil)
}
