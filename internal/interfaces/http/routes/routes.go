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
	salesHandler *handlers.SalesHandler) {

	// Public routes
	setupPublicRoutes(router, authHandler)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// User routes
	setupUserRoutes(api, userHandler)

	// Product routes
	setupProductRoutes(api, productHandler)

	// Sales routes
	setupSalesRoutes(api, salesHandler)
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
