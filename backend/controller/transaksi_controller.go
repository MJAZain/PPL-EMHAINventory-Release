package controller

import (
	"fmt"
	"go-gin-auth/dto"
	"go-gin-auth/model"
	"go-gin-auth/service"
	"go-gin-auth/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransaksiController struct {
	service service.TransaksiService
}

func NewTransaksiController(s service.TransaksiService) *TransaksiController {
	return &TransaksiController{service: s}
}

func (c *TransaksiController) CreateTransaksi(ctx *gin.Context) {
	var input dto.CreateTransaksiRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	// Mapping DTO to model
	transaksi := &model.Transaksi{
		NomorTransaksi:   model.GenerateNomorTransaksi(),
		ObatID:           input.ObatID,
		JumlahObat:       input.JumlahObat,
		TanggalPembelian: time.Now(),
		TotalHarga:       input.TotalHarga,
		UserID:           input.UserID,
	}
	transaksi, err := c.service.CreateTransaksi(*transaksi)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.Respond(ctx, http.StatusCreated, "Success", nil, transaksi)
}

func (c *TransaksiController) GetAllTransaksi(ctx *gin.Context) {
	transaksis, err := c.service.GetAllTransaksi()
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.Respond(ctx, http.StatusOK, "Success", nil, transaksis)
}
func (c *TransaksiController) DeleteTransaksi(ctx *gin.Context) {
	idParam := ctx.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, "error", "Invalid ID", nil)
		return
	}

	err := c.service.DeleteTransaksi(id)
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, "error", "Transaksi not found", nil)
		return
	}
	utils.Respond(ctx, http.StatusOK, "Transaksi deleted successfully", nil, nil)
}
