package router

import (
	"go-gin-auth/config"
	"go-gin-auth/controller"
	"go-gin-auth/internal/brand"
	"go-gin-auth/internal/category"
	"go-gin-auth/internal/doctor"
	"go-gin-auth/internal/drug_category"
	"go-gin-auth/internal/incomingProducts"
	"go-gin-auth/internal/location"
	"go-gin-auth/internal/nonpbf"
	"go-gin-auth/internal/outgoingProducts"
	"go-gin-auth/internal/patient"
	"go-gin-auth/internal/pbf"
	"go-gin-auth/internal/product"
	"go-gin-auth/internal/shift"
	"go-gin-auth/internal/stock_correction"
	storagelocation "go-gin-auth/internal/storage_location"
	"go-gin-auth/internal/supplier"
	"go-gin-auth/internal/unit"
	"go-gin-auth/middleware"
	"go-gin-auth/repository"
	"go-gin-auth/service"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() *gin.Engine {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()

	// Tambahkan ini
	// r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // alamat asal React kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	repo := repository.NewTransaksiRepository()
	svc := service.NewTransaksiService(repo)
	//ctrl := controller.NewTransaksiController(svc)
	api := r.Group("/api")
	{
		api.POST("/users/login", controller.Login)
		users := api.Group("/users")
		users.Use(middleware.AuthAdminMiddleware())
		{
			users.POST("/register", controller.Register)
			users.POST("/logout", controller.Logout)
			users.GET("/", controller.GetUsers)
			users.GET("/:id", controller.GetUser)
			users.PUT("/:id", controller.UpdateUser)
			users.DELETE("/:id", controller.DeleteUser)
			users.GET("/search", controller.SearchUsers)
			users.PATCH("/:id/deactivate", controller.DeactivateUser)
			users.PATCH("/:id/reactivate", controller.ReactivateUser)
			users.PUT("/:id/reset-password", controller.ResetUserPassword)
			users.GET("/export/csv", controller.ExportUsersCSV)

		}
		transaksi := api.Group("/transaksi")
		transaksi.Use(middleware.AuthAdminMiddleware()).DELETE("/:id", controller.NewTransaksiController(svc).DeleteTransaksi)
		transaksi.Use(middleware.AuthMiddleware())
		{
			transaksi.POST("/", controller.NewTransaksiController(svc).CreateTransaksi)
			transaksi.GET("/", controller.NewTransaksiController(svc).GetAllTransaksi)
		}
		unit := unit.NewUnitHandler()
		units := api.Group("/units")
		units.Use(middleware.AuthAdminMiddleware())
		{
			units.POST("/", unit.CreateUnit)
			units.GET("/", unit.GetUnits)
			units.GET("/:id", unit.GetUnitByID)
			units.PUT("/:id", unit.UpdateUnit)
			units.DELETE("/:id", unit.DeleteUnit)
		}

		category := category.NewCategoryHandler()
		categories := api.Group("/categories")
		categories.Use(middleware.AuthAdminMiddleware())
		{
			categories.POST("/", category.CreateCategory)
			categories.GET("/", category.GetCategories)
			categories.GET("/:id", category.GetCategoryByID)
			categories.PUT("/:id", category.UpdateCategory)
			categories.DELETE("/:id", category.DeleteCategory)
		}

		product := product.NewProductHandler()
		products := api.Group("/products")
		products.Use(middleware.AuthAdminMiddleware())
		{
			products.POST("/", product.CreateProduct)
			products.GET("/", product.GetProducts)
			products.GET("/:id", product.GetProductByID)
			products.PUT("/:id", product.UpdateProduct)
			products.DELETE("/:id", product.DeleteProduct)
		}

		// Stock Opname
		repoOpname := repository.NewStockOpnameRepository(config.DB)
		svcOpname := service.NewStockOpnameService(repoOpname)
		ctrlOpname := controller.NewStockOpnameController(svcOpname)

		// opname := api.Group("/opname")
		// {
		// 	opname.Use(middleware.AuthMiddleware()).POST("", ctrlOpname.Create)
		// 	opname.Use(middleware.AuthMiddleware()).GET("", ctrlOpname.GetAll)
		// 	opname.Use(middleware.AuthMiddleware()).GET("/:id", ctrlOpname.GetByID)
		// 	opname.Use(middleware.AuthAdminMiddleware()).DELETE("/:id", ctrlOpname.Delete)
		// }

		stockOpname := api.Group("/stock-opname")
		{
			//Reporting
			stockOpname.Use(middleware.AuthMiddleware()).GET("", ctrlOpname.GetOpnameList)
			stockOpname.Use(middleware.AuthMiddleware()).GET("/:opnameID", ctrlOpname.GetOpnameDetails)
			// Draft operations
			stockOpname.Use(middleware.AuthMiddleware()).POST("/draft", ctrlOpname.CreateDraft)
			stockOpname.Use(middleware.AuthMiddleware()).GET("/draft/:opnameID", ctrlOpname.GetDraft)
			stockOpname.Use(middleware.AuthMiddleware()).PUT("/draft/:opnameID", ctrlOpname.UpdateDraft)
			stockOpname.Use(middleware.AuthMiddleware()).DELETE("/draft/:opnameID", ctrlOpname.DeleteDraft)
			// Products operations
			stockOpname.Use(middleware.AuthMiddleware()).POST("/draft/:opnameID/products", ctrlOpname.AddProductToDraft)
			stockOpname.Use(middleware.AuthMiddleware()).DELETE("/draft/:opnameID/products/:detailID", ctrlOpname.RemoveProductFromDraft)
			// Process operations
			stockOpname.Use(middleware.AuthMiddleware()).POST("/:opnameID/start", ctrlOpname.StartOpname)
			stockOpname.Use(middleware.AuthMiddleware()).PUT("/details/:detailID/record", ctrlOpname.RecordActualStock)
			// Completion operations
			stockOpname.Use(middleware.AuthMiddleware()).POST("/:opnameID/complete", ctrlOpname.CompleteOpname)
			stockOpname.Use(middleware.AuthMiddleware()).POST("/:opnameID/cancel", ctrlOpname.CancelOpname)

			// users story
			stockOpname.Use(middleware.AuthMiddleware()).GET("/history", ctrlOpname.GetStockOpnameHistory)
			stockOpname.Use(middleware.AuthMiddleware()).GET("/products", ctrlOpname.GetProducts)
			stockOpname.Use(middleware.AuthMiddleware()).GET("/discrepancies", ctrlOpname.GetStockDiscrepancies)
			//otomatis tidak manual
			//stockOpname.Use(middleware.AuthMiddleware()).PUT("/products/:product_id", ctrlOpname.AdjustProductStock)
		}

		// Incoming Products
		incomingProduct := incomingProducts.NewHandlerIncomingProducts()
		incomingProductsGroup := api.Group("/incoming-products")
		incomingProductsGroup.Use(middleware.AuthAdminMiddleware())
		{
			incomingProductsGroup.POST("/", incomingProduct.CreateIncomingProduct)
			incomingProductsGroup.GET("/", incomingProduct.GetAllIncomingProducts)
			incomingProductsGroup.GET("/:id", incomingProduct.GetIncomingProductByID)
			incomingProductsGroup.PUT("/:id", incomingProduct.UpdateIncomingProduct)
			incomingProductsGroup.DELETE("/:id", incomingProduct.DeleteIncomingProduct)
		}

		outgoingProduct := outgoingProducts.NewHandlerOutgoingProducts()
		outgoingProductGroup := api.Group("/outgoing-products")
		outgoingProductGroup.Use(middleware.AuthAdminMiddleware())
		{
			outgoingProductGroup.POST("/", outgoingProduct.CreateOutgoingProduct)
			outgoingProductGroup.GET("/", outgoingProduct.GetAllOutgoingProducts)
			outgoingProductGroup.GET("/:id", outgoingProduct.GetOutgoingProductByID)
			outgoingProductGroup.PUT("/:id", outgoingProduct.UpdateOutgoingProduct)
			outgoingProductGroup.DELETE("/:id", outgoingProduct.DeleteOutgoingProduct)
		}

		apiAuth := api
		apiAuth.Use(middleware.AuthAdminMiddleware())
		storagelocation.StorageLocationRouter(apiAuth)
		brand.BrandRouter(apiAuth)
		supplier.SupplierRouter(apiAuth)
		location.LocationRouter(apiAuth)
		doctor.DoctorRouter(apiAuth)
		patient.PatientRouter(apiAuth)
		drug_category.DrugCategoryRouter(apiAuth)
		shift.ShiftRouter(apiAuth)
		stock_correction.StockCorrectionRouter(apiAuth)

		pbfRouter := api.Group("/incoming-pbf")
		pbfRouter.Use(middleware.AuthMiddleware()).GET("", pbf.GetAllIncomingPBF)
		pbfRouter.Use(middleware.AuthMiddleware()).POST("", pbf.CreateIncomingPBF)
		pbfRouter.Use(middleware.AuthMiddleware()).GET("/:id", pbf.GetIncomingPBFByID)
		pbfRouter.Use(middleware.AuthMiddleware()).PUT("/:id", pbf.UpdateIncomingPBF)
		pbfRouter.Use(middleware.AuthMiddleware()).DELETE("/:id", pbf.DeleteIncomingPBF)

		nonpbfService := nonpbf.NewIncomingNonPBFService(config.DB)
		nonpbfController := nonpbf.NewIncomingNonPBFController(nonpbfService)

		nonpbfRouter := api.Group("/incoming-nonpbf", middleware.AuthMiddleware())
		nonpbfRouter.GET("", nonpbfController.GetAll)
		nonpbfRouter.POST("", nonpbfController.Create)
		nonpbfRouter.GET("/:id", nonpbfController.GetByID)
		nonpbfRouter.PUT("/:id", nonpbfController.Update)
		nonpbfRouter.DELETE("/:id", nonpbfController.Delete)

	}
	return r
}
