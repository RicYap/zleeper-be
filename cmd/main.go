package main

import (
	"log"

	"zleeper-be/config"
	"zleeper-be/internal/controllers"
	repositories "zleeper-be/internal/datas"
	"zleeper-be/internal/routes"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"
	"zleeper-be/pkg/cache"
	"zleeper-be/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger("app.log")

	// Initialize database
	db, err := database.InitDB(cfg.DBConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient := cache.NewRedisClient(cfg.RedisConfig)
	defer redisClient.Close()

	// Auto migrate models
	if err := database.MigrateDB(db); err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	e.Use(middleware.Recover())
	e.Use(utils.RequestLoggerMiddleware(logger))

	// Initialize repositories
	orderItemRepo := repositories.NewOrderItemData(db)
	userRepo := repositories.NewUserData(db)
	orderHistoryRepo := repositories.NewOrderHistoryData(db)

	// Initialize services
	orderItemService := services.NewOrderItemService(orderItemRepo, redisClient)
	userService := services.NewUserService(userRepo, redisClient)
	orderHistoryService := services.NewOrderHistoryService(orderHistoryRepo, redisClient)

	// Initialize controllers
	orderItemController := controllers.NewOrderItemController(orderItemService)
	userController := controllers.NewUserController(userService)
	orderHistoryController := controllers.NewOrderHistoryController(orderHistoryService)

	// Register routes
	routes.NewRoutes(
		orderItemController,
		userController,
		orderHistoryController,
	).RegisterRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}