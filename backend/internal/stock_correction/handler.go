package stock_correction

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

func (h *Handler) CreateCorrection(c *gin.Context) {
	var input StockCorrection
	if err := c.ShouldBind(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Input tidak valid", err.Error(), nil)
		return
	}

	officerName, _ := c.Get("full_name")

	newCorrection, err := h.service.CreateCorrection(&input, officerName.(string))
	if err != nil {
		if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrStockUpdateFailed) {
			utils.Respond(c, http.StatusBadRequest, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal membuat koreksi stok", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusCreated, "Koreksi stok berhasil dibuat", nil, newCorrection)
}

func (h *Handler) GetAllCorrections(c *gin.Context) {
	corrections, err := h.service.GetAllCorrections()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil histori koreksi", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Histori koreksi berhasil diambil", nil, corrections)
}

func (h *Handler) GetCorrectionByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	correction, err := h.service.GetCorrectionByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data koreksi", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Detail koreksi berhasil diambil", nil, correction)
}

func (h *Handler) DeleteCorrection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteCorrection(uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Respond(c, http.StatusNotFound, err.Error(), err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "Gagal menghapus koreksi", err.Error(), nil)
		}
		return
	}
	utils.Respond(c, http.StatusOK, "Koreksi berhasil dihapus", nil, nil)
}
