package expense

import (
	"go-gin-auth/config"
	"go-gin-auth/internal/expense_type"
	"go-gin-auth/middleware"

	"github.com/gin-gonic/gin"
)

func ExpenseRouter(api *gin.RouterGroup) {
	expenseTypeRepo := expense_type.NewRepository(config.DB)
	expenseTypeService := expense_type.NewService(expenseTypeRepo)

	repo := NewRepository(config.DB)
	service := NewService(repo, expenseTypeService)
	handler := NewHandler(service)

	expenseGroup := api.Group("/expenses")
	expenseGroup.Use(middleware.AuthMiddleware())
	{
		expenseGroup.POST("/", handler.CreateExpense)
		expenseGroup.GET("/", handler.GetAllExpenses)
		expenseGroup.GET("/:id", handler.GetExpenseByID)
		expenseGroup.PUT("/:id", handler.UpdateExpense)
		expenseGroup.DELETE("/:id", handler.DeleteExpense)
	}
}
