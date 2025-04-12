package main

import (
	"log"
	"os"

	"Sheikh-Enterprise-Backend/internal/config"
	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/middleware"
	"Sheikh-Enterprise-Backend/pkg/database"
	"Sheikh-Enterprise-Backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Initialize(&cfg.Logger); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize validator
	validator.Initialize()

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("ALLOWED_ORIGINS")}
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	// Add global middleware
	router.Use(middleware.ValidationMiddleware())

	// Initialize routes
	setupRoutes(router)

	// Start server
	logger.Info("Starting server on port " + cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func setupRoutes(router *gin.Engine) {
	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
