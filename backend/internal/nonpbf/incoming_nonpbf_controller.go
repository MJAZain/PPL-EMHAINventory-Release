package nonpbf

import (
	"fmt"
	"go-gin-auth/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IncomingNonPBFController struct {
	service IncomingNonPBFServiceInterface
}

func NewIncomingNonPBFController(service IncomingNonPBFServiceInterface) *IncomingNonPBFController {
	return &IncomingNonPBFController{service: service}
}

func (ctrl *IncomingNonPBFController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	incomings, total, err := ctrl.service.GetAll(page, limit)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data", err.Error(), nil)
		return
	}

	response := map[string]interface{}{
		"data":        incomings,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}

	utils.Respond(c, http.StatusOK, "Data berhasil diambil", nil, response)
}

func (ctrl *IncomingNonPBFController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "ID tidak valid", err.Error(), nil)
		return
	}

	incoming, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			utils.Respond(c, http.StatusNotFound, "Data tidak ditemukan", err.Error(), nil)
			return
		}
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengambil data", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data berhasil diambil", nil, incoming)
}

func (ctrl *IncomingNonPBFController) Create(c *gin.Context) {
	var req CreateIncomingNonPBFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Data tidak valid", err.Error(), nil)
		return
	}

	incoming, err := ctrl.service.Create(req)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Gagal menyimpan data", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusCreated, "Data berhasil disimpan", nil, incoming)
}

func (ctrl *IncomingNonPBFController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "ID tidak valid", err.Error(), nil)
		return
	}

	var req UpdateIncomingNonPBFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Data tidak valid", err.Error(), nil)
		return
	}

	incoming, err := ctrl.service.Update(uint(id), req)
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			utils.Respond(c, http.StatusNotFound, "Data tidak ditemukan", err.Error(), nil)
			return
		}
		utils.Respond(c, http.StatusInternalServerError, "Gagal mengupdate data", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data berhasil diupdate", nil, incoming)
}

func (ctrl *IncomingNonPBFController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	fmt.Println("ID:", id)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "ID tidak valid", err.Error(), nil)
		return
	}

	err = ctrl.service.Delete(uint(id))
	if err != nil {
		if err.Error() == "data tidak ditemukan" {
			utils.Respond(c, http.StatusNotFound, "Data tidak ditemukan", err.Error(), nil)
			return
		}
		utils.Respond(c, http.StatusInternalServerError, "Gagal menghapus data", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Data berhasil dihapus", nil, nil)
}
