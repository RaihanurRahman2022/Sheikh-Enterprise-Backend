package routes

import (
	"Sheikh-Enterprise-Backend/internal/interfaces/http/handlers"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	salesHandler *handlers.SalesHandler,
	supplierHandler *handlers.SupplierHandler,
	purchaseHandler *handlers.PurchaseHandler) {

	// Public routes
	setupPublicRoutes(router, authHandler)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// User routes
	setupUserRoutes(api, userHandler)

	// Master Data routes
	setupProductRoutes(api, productHandler)
	setupCompanyRoutes(api, companyHandler)
	setupShopRoutes(api, shopHandler)
	setupCustomerRoutes(api, customerHandler)
	setupSupplierRoutes(api, supplierHandler)

	// Transaction routes
	setupSalesRoutes(api, salesHandler)
	setupPurchaseRoutes(api, purchaseHandler)
	setupStockTransferRoutes(api, stockTransferHandler)
	setupPaymentRoutes(api, paymentHandler)

	// Analytics routes
	setupAnalyticsRoutes(api, analyticsHandler)
}

// setupSupplierRoutes configures supplier-related routes
func setupSupplierRoutes(api *gin.RouterGroup, supplierHandler *handlers.SupplierHandler) {
	suppliers := api.Group("/suppliers")
	{
		suppliers.GET("", supplierHandler.GetSuppliers)
		suppliers.GET("/:id", supplierHandler.GetSupplier)
		suppliers.POST("", supplierHandler.CreateSupplier)
		suppliers.DELETE("/:id", supplierHandler.DeleteSupplier)
	}
}

// setupPurchaseRoutes configures purchase-related routes
func setupPurchaseRoutes(api *gin.RouterGroup, purchaseHandler *handlers.PurchaseHandler) {
	purchases := api.Group("/purchases")
	{
		purchases.GET("", purchaseHandler.GetPurchases)
		purchases.GET("/:id", purchaseHandler.GetPurchase)
		purchases.POST("", purchaseHandler.CreatePurchase)
		purchases.DELETE("/:id", purchaseHandler.DeletePurchase)
	}
}

// setupPublicRoutes configures public routes that don't require authentication
func setupPublicRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		users.GET("/me", userHandler.GetUserDetails)
		users.PUT("/me", userHandler.UpdateUserDetails)
		users.PUT("/change-password", userHandler.UpdatePassword)
	}
}

// setupProductRoutes configures product-related routes
func setupProductRoutes(api *gin.RouterGroup, productHandler *handlers.ProductHandler) {
	products := api.Group("/products")
	{
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.POST("", productHandler.CreateProduct)
		products.GET("/export", productHandler.ExportToExcel)
		products.POST("/bulk-import", productHandler.BulkImport)
	}
}

// setupSalesRoutes configures sales-related routes
func setupSalesRoutes(api *gin.RouterGroup, salesHandler *handlers.SalesHandler) {
	sales := api.Group("/sales")
	{
		sales.GET("", salesHandler.GetSales)
		sales.GET("/:id", salesHandler.GetSale)
		sales.POST("", salesHandler.CreateSale)
		sales.DELETE("/:id", salesHandler.DeleteSale)
		sales.GET("/export", salesHandler.ExportToExcel)

		// Analytics routes
		analytics := sales.Group("/analytics")
		{
			analytics.GET("", salesHandler.GetAnalytics)
			analytics.GET("/last-7-days", salesHandler.GetLast7DaysSales)
		}
	}
}
