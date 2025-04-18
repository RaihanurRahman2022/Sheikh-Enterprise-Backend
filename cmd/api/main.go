package main

import (
	"log"
	"os"
	"strings"
	"time" // Add this import for cors.MaxAge

	"Sheikh-Enterprise-Backend/internal/config"
	repository "Sheikh-Enterprise-Backend/internal/infrastructure/persistence"
	validator "Sheikh-Enterprise-Backend/internal/infrastructure/validation"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/handlers"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/middleware"
	"Sheikh-Enterprise-Backend/internal/interfaces/http/routes"
	services "Sheikh-Enterprise-Backend/internal/usecases/impl"
	"Sheikh-Enterprise-Backend/pkg/database"
	"Sheikh-Enterprise-Backend/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Initialize application
	app, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start server
	logger.Info("Starting server on port " + app.Config.Server.Port)
	if err := app.Router.Run(":" + app.Config.Server.Port); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

// App represents the main application structure
type App struct {
	Config   *config.Config
	Router   *gin.Engine
	Handlers *handlers.Handlers
}

// initializeApp sets up the application components
func initializeApp() (*App, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	if err := logger.Initialize(&cfg.Logger); err != nil {
		return nil, err
	}

	// Initialize validator
	validator.Initialize()

	// Initialize database
	db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	repos := initializeRepositories(db)

	// Initialize services
	svcs := initializeServices(repos)

	// Initialize handlers
	handlers := initializeHandlers(svcs)

	// Create and configure router
	router := configureRouter(handlers)

	return &App{
		Config:   cfg,
		Router:   router,
		Handlers: handlers,
	}, nil
}

// initializeDatabase sets up the database connection and runs migrations
func initializeDatabase() (*gorm.DB, error) {
	db, err := database.InitDB()
	if err != nil {
		return nil, err
	}

	if err := database.RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// initializeRepositories creates all repository instances
func initializeRepositories(db *gorm.DB) *repository.Repositories {
	return &repository.Repositories{
		User:     repository.NewUserRepository(db),
		Product:  repository.NewProductRepository(db),
		Sales:    repository.NewSalesRepository(db),
		Purchase: repository.NewPurchaseRepository(db),
		Supplier: repository.NewSupplierRepository(db),
		Company:  repository.NewCompanyRepository(db),
		Shop:     repository.NewShopRepository(db),
	}
}

// initializeServices creates all service instances
func initializeServices(repos *repository.Repositories) *services.Services {
	return &services.Services{
		Auth:     services.NewAuthService(repos.Auth),
		User:     services.NewUserService(repos.User),
		Product:  services.NewProductService(repos.Product),
		Sales:    services.NewSalesService(repos.Sales),
		Purchase: services.NewPurchaseService(repos.Purchase),
		Supplier: services.NewSupplierService(repos.Supplier),
		Company:  services.NewCompanyService(repos.Company),
		Shop:     services.NewShopService(repos.Shop),
	}
}

// initializeHandlers creates all handler instances
func initializeHandlers(svcs *services.Services) *handlers.Handlers {
	return &handlers.Handlers{
		Auth:     handlers.NewAuthHandler(svcs.Auth),
		User:     handlers.NewUserHandler(svcs.User),
		Product:  handlers.NewProductHandler(svcs.Product),
		Sales:    handlers.NewSalesHandler(svcs.Sales),
		Purchase: handlers.NewPurchaseHandler(svcs.Purchase),
		Supplier: handlers.NewSupplierHandler(svcs.Supplier),
		Company:  handlers.NewCompanyHandler(svcs.Company),
		Shop:     handlers.NewShopHandler(svcs.Shop),
	}
}

// configureRouter sets up the Gin router with middleware and routes
func configureRouter(handlers *handlers.Handlers) *gin.Engine {
	router := gin.Default()

	// Set trusted proxies (none for development)
	router.SetTrustedProxies(nil) // Or []string{"127.0.0.1"} for localhost

	// Configure CORS
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}
	origins := strings.Split(allowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	corsConfig := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ","),
		AllowHeaders:     strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ","),
		ExposeHeaders:    strings.Split(os.Getenv("CORS_EXPOSE_HEADERS"), ","),
		AllowCredentials: true,
		MaxAge:           300 * time.Second,
	}

	router.Use(cors.New(corsConfig))

	// Add global middleware
	router.Use(middleware.ValidationMiddleware())

	// Setup routes
	routes.SetupRoutes(router, handlers)

	return router
}
