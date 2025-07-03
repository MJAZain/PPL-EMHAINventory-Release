package location

import "github.com/gin-gonic/gin"

func LocationRouter(api *gin.RouterGroup) Service {
	locationService := NewService()

	locationHandler := NewHandler(locationService)

	locationGroup := api.Group("/locations")
	{
		locationGroup.GET("/provinces", locationHandler.GetProvinces)

		locationGroup.GET("/regencies/:province_id", locationHandler.GetRegenciesByProvinceID)
	}

	return locationService
}
