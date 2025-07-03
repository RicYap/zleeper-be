package routes

import "github.com/labstack/echo/v4"

func (r *Routes) registerOrderItemRoutes(api *echo.Group) {
	orderItems := api.Group("/order-items")
	orderItems.POST("", r.OrderItemController.Create)
	orderItems.GET("", r.OrderItemController.List)
	orderItems.GET("/:id", r.OrderItemController.Get)
	orderItems.PUT("/:id", r.OrderItemController.Update)
	orderItems.DELETE("/:id", r.OrderItemController.Delete)
}