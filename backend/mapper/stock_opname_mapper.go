package mapper

import (
	"go-gin-auth/dto"
	"go-gin-auth/internal/opname"
	"strconv"
)

// Fungsi untuk mengonversi StockOpnameRequest menjadi StockOpname
func ToModelStockOpname(request dto.StockOpnameRequest) opname.StockOpname {
	var details []opname.StockOpnameDetail
	// Pemetaan untuk detail
	for _, detail := range request.Details {
		details = append(details, opname.StockOpnameDetail{
			ProductID:   detail.ObatID,
			ActualStock: detail.StokFisik,
			// StokSistem: 0, // Asumsikan nilai stok sistem untuk detail, sesuaikan jika ada data lain
			// Selisih:    0, // Asumsikan nilai selisih, sesuaikan sesuai dengan kebutuhan
		})
	}

	// Pemetaan objek StockOpname utama
	return opname.StockOpname{
		CreatedBy: strconv.FormatUint(uint64(request.UserID), 10),
		Details:   details,
	}
}
