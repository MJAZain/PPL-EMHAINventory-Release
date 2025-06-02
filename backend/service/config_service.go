package service

import (
	"go-gin-auth/config"
	"go-gin-auth/model"
)

func GetLoginConfig() (model.SystemConfig, error) {
	var data model.SystemConfig
	err := config.DB.First(&data).Error
	return data, err
}

func UpdateLoginConfig(updatedConfig model.SystemConfig) error {
	return config.DB.Save(&updatedConfig).Error
}
