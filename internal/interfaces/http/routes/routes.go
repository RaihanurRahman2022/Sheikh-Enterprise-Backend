package routes

import (
	"Sheikh-Enterprise-Backend/internal/interfaces/http/handlers"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *gin.Engine, handlers *handlers.Handlers) {
	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Sheikh Enterprise Inventory Management System API",
			"version": "1.0.0",
			"endpoints": []string{
				"/health",
				"/auth/login",
				"/api/shops",
				"/api/companies",
				"/api/products",
			},
		})
	})

	// Public routes (no auth required)
	setupPublicRoutes(router, handlers.Auth)

	// API routes (auth required)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		setupUserRoutes(api, handlers.User)
		setupProductRoutes(api, handlers.Product)
		setupSalesRoutes(api, handlers.Sales)
		setupPurchaseRoutes(api, handlers.Purchase)
		setupSupplierRoutes(api, handlers.Supplier)
		setupCompanyRoutes(api, handlers.Company)
		setupShopRoutes(api, handlers.Shop)
	}
}

// setupPublicRoutes configures public routes that don't require authentication
func setupPublicRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/change-password", authHandler.ChangePassword)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}

// setupUserRoutes configures user-related routes
func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		users.GET("/me", userHandler.GetUserDetails)
		users.PUT("/me", userHandler.UpdateUserDetails)
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

// setupSupplierRoutes configures supplier-related routes
func setupSupplierRoutes(api *gin.RouterGroup, supplierHandler *handlers.SupplierHandler) {
	suppliers := api.Group("/suppliers")
	{
		suppliers.GET("", supplierHandler.GetSuppliers)
		suppliers.GET("/:id", supplierHandler.GetSupplier)
		suppliers.POST("", supplierHandler.CreateSupplier)
		suppliers.PUT("/:id", supplierHandler.UpdateSupplier)
		suppliers.DELETE("/:id", supplierHandler.DeleteSupplier)
	}
}

// setupCompanyRoutes configures company-related routes
func setupCompanyRoutes(api *gin.RouterGroup, companyHandler *handlers.CompanyHandler) {
	companies := api.Group("/companies")
	{
		companies.GET("", companyHandler.GetCompanies)
		companies.GET("/:id", companyHandler.GetCompany)
		companies.POST("", companyHandler.CreateCompany)
		companies.PUT("/:id", companyHandler.UpdateCompany)
		companies.DELETE("/:id", companyHandler.DeleteCompany)
	}
}

// setupShopRoutes configures shop-related routes
func setupShopRoutes(api *gin.RouterGroup, shopHandler *handlers.ShopHandler) {
	shops := api.Group("/shops")
	{
		shops.GET("", shopHandler.GetShops)
		shops.GET("/:id", shopHandler.GetShop)
		shops.POST("", shopHandler.CreateShop)
		shops.PUT("/:id", shopHandler.UpdateShop)
		shops.DELETE("/:id", shopHandler.DeleteShop)
		shops.GET("/company/:company_id", shopHandler.GetShopsByCompany)
	}
}
