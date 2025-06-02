package service

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
	"time"

	"github.com/gin-gonic/gin"
)

// LogActivity mencatat aktivitas pengguna
func LogActivity(userID uint, fullName string, activityType string, description string, c *gin.Context) error {
	ipAddress := c.ClientIP()          // Ambil IP Address
	userAgent := c.Request.UserAgent() // Ambil User-Agent

	// Membuat objek log aktivitas
	log := model.ActivityLog{
		UserID:       userID,
		FullName:     fullName,
		ActivityType: activityType,
		Description:  description,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ActivityTime: time.Now(),
	}

	// Menyimpan log ke database
	if err := config.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// GetActivityLogs mengambil log aktivitas berdasarkan filter tertentu
func GetActivityLogs(filters map[string]string) ([]model.ActivityLog, error) {
	var logs []model.ActivityLog
	query := config.DB.Model(&model.ActivityLog{})

	// Filter berdasarkan username jika ada
	if filters["username"] != "" {
		query = query.Where("username LIKE ?", "%"+filters["username"]+"%")
	}

	// Filter berdasarkan jenis aktivitas jika ada
	if filters["activity_type"] != "" {
		query = query.Where("activity_type LIKE ?", "%"+filters["activity_type"]+"%")
	}

	// Ambil log aktivitas
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
