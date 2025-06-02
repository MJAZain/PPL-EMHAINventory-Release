package brand

import "github.com/gin-gonic/gin"

func BrandRouter(api *gin.RouterGroup) {
	brandRepo := NewBrandRepository()
	brandService := NewBrandService(brandRepo)
	brandHandler := NewBrandHandler(brandService)

	brandsAPI := api.Group("/brands")
	{
		brandsAPI.POST("", brandHandler.CreateBrand)
		brandsAPI.GET("", brandHandler.GetBrands)
		brandsAPI.GET("/:id", brandHandler.GetBrandByID)
		brandsAPI.PUT("/:id", brandHandler.UpdateBrand)
		brandsAPI.DELETE("/:id", brandHandler.DeleteBrand)
	}
}
