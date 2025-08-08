package analysis

import (
	"go-gin-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

// NewHandler creates and returns a new Handler with the provided Service.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAnalysisData handles analysis data request
// @Summary Get analysis data
// @Description Get data for analysis
// @Tags Analysis
// @Accept json
// @Produce json
// @Param start_date query string false "Start date for analysis"
// @Param end_date query string false "End date for analysis"
// @Success 200 {object} []AnalysisResult
// @Failure 400 {object} ErrorResponse
// @Router /api/analysis [get]
func (h *Handler) GetAnalysisData(c *gin.Context) {
	var req AnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Parameter query tidak valid", err.Error(), nil)
		return
	}

	if (req.StartDate != "" && req.EndDate == "") || (req.StartDate == "" && req.EndDate != "") {
		utils.Respond(c, http.StatusBadRequest, "Parameter tidak lengkap", "Gunakan start_date dan end_date secara bersamaan.", nil)
		return
	}

	data, err := h.service.GetAnalysis(req)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Gagal memproses analisis", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data analisis berhasil diambil", nil, data)
}
