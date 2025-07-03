package supplier

import (
	"go-gin-auth/config"
	"go-gin-auth/internal/location"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func SupplierRouter(api *gin.RouterGroup) {
	supplierRepo := NewRepository(config.DB)

	locationService := location.NewService()
	supplierService := NewService(supplierRepo, locationService)

	supplierHandler := NewHandler(supplierService)

	supplierGroup := api.Group("/suppliers")
	supplierGroup.Use(middleware.AuthAdminMiddleware())
	{
		supplierGroup.POST("/", supplierHandler.CreateSupplier)
		supplierGroup.GET("/", supplierHandler.GetAllSuppliers)
		supplierGroup.GET("/:id", supplierHandler.GetSupplierByID)
		supplierGroup.PUT("/:id", supplierHandler.UpdateSupplier)
		supplierGroup.DELETE("/:id", supplierHandler.DeleteSupplier)
	}
}
