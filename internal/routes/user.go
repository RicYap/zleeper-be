package routes

import "github.com/labstack/echo/v4"

func (r *Routes) registerUserRoutes(api *echo.Group) {
	users := api.Group("/users")
	users.POST("", r.UserController.Create)
	users.GET("", r.UserController.List)
	users.GET("/:id", r.UserController.Get)
	users.PUT("/:id", r.UserController.Update)
	users.DELETE("/:id", r.UserController.Delete)
}