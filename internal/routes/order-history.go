package routes

import "github.com/labstack/echo/v4"

func (r *Routes) registerOrderHistoryRoutes(api *echo.Group) {
	orderItems := api.Group("/order-histories")
	orderItems.POST("", r.OrderHistoryController.Create)
	orderItems.GET("", r.OrderHistoryController.List)
	orderItems.GET("/:id", r.OrderHistoryController.Get)
	orderItems.PUT("/:id", r.OrderHistoryController.Update)
	orderItems.DELETE("/:id", r.OrderHistoryController.Delete)
}