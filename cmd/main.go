package main

import (
	"log"

	"zleeper-be/config"
	"zleeper-be/internal/controllers"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/routes"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"
	"zleeper-be/pkg/cache"
	"zleeper-be/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := utils.NewLogger("app.log")

	db, err := database.InitDB(cfg.DBConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := cache.NewRedisClient(cfg.RedisConfig)
	defer redisClient.Close()

	if err := database.MigrateDB(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: logger.Writer(),
	}))
	e.Use(middleware.Recover())
	e.Use(utils.RequestLoggerMiddleware(logger))

	orderItemData := datas.NewOrderItemData(db)
	userData := datas.NewUserData(db)
	orderHistoryData := datas.NewOrderHistoryData(db)

	orderItemService := services.NewOrderItemService(orderItemData, redisClient)
	userService := services.NewUserService(userData, redisClient)
	orderHistoryService := services.NewOrderHistoryService(orderHistoryData, redisClient)

	orderItemController := controllers.NewOrderItemController(orderItemService)
	userController := controllers.NewUserController(userService)
	orderHistoryController := controllers.NewOrderHistoryController(orderHistoryService)

	routes.NewRoutes(
		orderItemController,
		userController,
		orderHistoryController,
	).RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}