package analysis

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func AnalysisRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	analysisGroup := api.Group("/analysis")
	analysisGroup.Use(middleware.AuthMiddleware())
	{
		analysisGroup.GET("/", handler.GetAnalysisData)
	}
}
