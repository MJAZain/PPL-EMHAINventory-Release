package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

func Respond(c *gin.Context, status int, message string, err interface{}, data interface{}) {
	c.JSON(status, APIResponse{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
	})
}
