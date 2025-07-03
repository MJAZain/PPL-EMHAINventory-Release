package patient

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func PatientRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	patientGroup := api.Group("/patients")
	patientGroup.Use(middleware.AuthAdminMiddleware())
	{
		patientGroup.POST("/", handler.CreatePatient)
		patientGroup.GET("/", handler.GetAllPatients)
		patientGroup.GET("/:id", handler.GetPatientByID)
		patientGroup.PUT("/:id", handler.UpdatePatient)
		patientGroup.DELETE("/:id", handler.DeletePatient)
	}
}
