package location

import (
	"go-gin-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProvinces(c *gin.Context) {
	provinces, err := h.service.GetAllProvinces()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data provinsi", err, nil)
		return
	}
	utils.Respond(c, http.StatusOK, "Data provinsi berhasil diambil", nil, provinces)
}

func (h *Handler) GetRegenciesByProvinceID(c *gin.Context) {
	provinceID := c.Param("province_id")
	if provinceID == "" {
		utils.Respond(c, http.StatusBadRequest, "Parameter province_id dibutuhkan", nil, nil)
		return
	}

	regencies, err := h.service.GetRegenciesByProvinceID(provinceID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data kota/kabupaten", err, nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data kota/kabupaten berhasil diambil", nil, regencies)
}
