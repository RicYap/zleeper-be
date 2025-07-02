package routes

import (
	"zleeper-be/internal/controllers"

	"github.com/labstack/echo/v4"
)

type Routes struct {
	OrderItemController     *controllers.OrderItemController
	UserController         *controllers.UserController
	OrderHistoryController *controllers.OrderHistoryController
}

func NewRoutes(
	orderItemController *controllers.OrderItemController,
	userController *controllers.UserController,
	orderHistoryController *controllers.OrderHistoryController,
) *Routes {
	return &Routes{
		OrderItemController:     orderItemController,
		UserController:         userController,
		OrderHistoryController: orderHistoryController,
	}
}

func (r *Routes) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api")

	// Register all route groups
	r.registerOrderItemRoutes(api)
	r.registerUserRoutes(api)
	r.registerOrderHistoryRoutes(api)
}