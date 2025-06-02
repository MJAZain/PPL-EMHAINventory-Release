package storagelocation

import "github.com/gin-gonic/gin"

func StorageLocationRouter(api *gin.RouterGroup) {
	storageLocationRepo := NewStorageLocationRepository()
	storageLocationService := NewStorageLocationService(storageLocationRepo)
	storageLocationHandler := NewStorageLocationHandler(storageLocationService)

	storageLocationsAPI := api.Group("/storage-locations")
	{
		storageLocationsAPI.POST("", storageLocationHandler.CreateStorageLocation)
		storageLocationsAPI.GET("", storageLocationHandler.GetStorageLocations)
		storageLocationsAPI.GET("/:id", storageLocationHandler.GetStorageLocationByID)
		storageLocationsAPI.PUT("/:id", storageLocationHandler.UpdateStorageLocation)
		storageLocationsAPI.DELETE("/:id", storageLocationHandler.DeleteStorageLocation)
	}
}
