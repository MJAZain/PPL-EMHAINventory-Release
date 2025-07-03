package stock_correction

import (
	"go-gin-auth/config"
	"go-gin-auth/internal/stock"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func StockCorrectionRouter(api *gin.RouterGroup) {
	stockRepo := stock.NewRepository()
	repo := NewRepository(config.DB)
	service := NewService(repo, stockRepo)
	handler := NewHandler(service)

	correctionGroup := api.Group("/stock-corrections")
	correctionGroup.Use(middleware.AuthAdminMiddleware())
	{
		correctionGroup.POST("/", handler.CreateCorrection)
		correctionGroup.GET("/", handler.GetAllCorrections)
		correctionGroup.GET("/:id", handler.GetCorrectionByID)
		correctionGroup.DELETE("/:id", handler.DeleteCorrection)
	}
}
