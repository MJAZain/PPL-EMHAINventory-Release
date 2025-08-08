package expense_type

import (
	"go-gin-auth/config"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func ExpenseTypeRouter(api *gin.RouterGroup) {
	repo := NewRepository(config.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	expenseTypeGroup := api.Group("/expense-types")
	expenseTypeGroup.Use(middleware.AuthAdminMiddleware())
	{
		expenseTypeGroup.POST("/", handler.CreateExpenseType)
		expenseTypeGroup.GET("/", handler.GetAllExpenseTypes)
		expenseTypeGroup.GET("/:id", handler.GetExpenseTypeByID)
		expenseTypeGroup.PUT("/:id", handler.UpdateExpenseType)
		expenseTypeGroup.DELETE("/:id", handler.DeleteExpenseType)
	}
}
