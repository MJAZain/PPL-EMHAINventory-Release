package utils

import (
	"go-gin-auth/dto"
	"go-gin-auth/model"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Ambil user_id dari context dan ubah jadi uint
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	// Asumsi user_id di JWT disimpan sebagai float64 (default)
	if idFloat, ok := userID.(float64); ok {
		return uint(idFloat)
	}

	return 0
}
func GetTableName(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Ambil tipe yang ditunjuk oleh pointer
	}
	// Ambil nama struct
	tableName := t.Name()
	return strings.ToLower(tableName) // Menyesuaikan dengan konvensi, misalnya nama tabel kecil
}
func ConvertDTOToUser(userDTO dto.RegisterRequestDTO) (*model.User, error) {
	var user model.User

	// Pemetaan otomatis antara DTO dan model menggunakan mapstructure
	err := mapstructure.Decode(userDTO, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
