package drug_category

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func DrugCategoryRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	categoryGroup := api.Group("/drug-categories")
	categoryGroup.Use(middleware.AuthAdminMiddleware())
	{
		categoryGroup.POST("/", handler.CreateCategory)
		categoryGroup.GET("/", handler.GetAllCategories)
		categoryGroup.GET("/:id", handler.GetCategoryByID)
		categoryGroup.PUT("/:id", handler.UpdateCategory)
		categoryGroup.DELETE("/:id", handler.DeleteCategory)
	}
}
