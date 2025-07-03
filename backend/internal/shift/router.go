package shift

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func ShiftRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	shiftGroup := api.Group("/shifts")
	shiftGroup.Use(middleware.AuthAdminMiddleware())
	{
		shiftGroup.POST("/open", handler.OpenShift)
		shiftGroup.PUT("/close/:id", handler.CloseShift)

		shiftGroup.GET("/", handler.GetAllShifts)
		shiftGroup.GET("/:id", handler.GetShiftByID)
		shiftGroup.PUT("/:id", handler.UpdateShift)
		shiftGroup.DELETE("/:id", handler.DeleteShift)
	}
}
