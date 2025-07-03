package doctor

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func DoctorRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	doctorGroup := api.Group("/doctors")
	doctorGroup.Use(middleware.AuthAdminMiddleware())
	{
		doctorGroup.POST("/", handler.CreateDoctor)
		doctorGroup.GET("/", handler.GetAllDoctors)
		doctorGroup.GET("/:id", handler.GetDoctorByID)
		doctorGroup.PUT("/:id", handler.UpdateDoctor)
		doctorGroup.DELETE("/:id", handler.DeleteDoctor)
	}
}
