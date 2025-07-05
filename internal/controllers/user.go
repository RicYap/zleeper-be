package controllers

import (
	"net/http"
	"strconv"
	"zleeper-be/internal/models"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service services.UserService
}

func (c *UserController) Create(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.service.Create(ctx.Request().Context(), &user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user")
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, user)
}

func (c *UserController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	orderItems, err := c.service.List(ctx.Request().Context(), page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch users")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, orderItems)
}

func (c *UserController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	user, err := c.service.Get(ctx.Request().Context(), uint(id))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "User not found")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, user)
}

func (c *UserController) Update(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	user.ID = uint(id)
	if err := c.service.Update(ctx.Request().Context(), &user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update user")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, user)
}

func (c *UserController) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	if err := c.service.Delete(ctx.Request().Context(), uint(id)); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete user")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, map[string]interface{}{"message": "User deleted successfully"})
}