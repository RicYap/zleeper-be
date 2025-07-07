package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"zleeper-be/internal/models"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserService
}

func (c *UserController) Create(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	resultChan := make(chan uint)
	errChan := make(chan error)

	go func() {
		if err := c.service.Create(ctx.Request().Context(), &user); err != nil {
			errChan <- err
			return
		}
		resultChan <- user.ID
	}()

	select {
	case id := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusCreated, map[string]interface{}{
			"id": id,
		})
	case err := <-errChan:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user: "+err.Error())
	}
}

func (c *UserController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resultChan := make(chan models.UserPagination)
	errChan := make(chan error)

	go func() {
		result, err := c.service.List(ctx.Request().Context(), page, limit)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return utils.PaginatedResponse(ctx, http.StatusOK, result.Data, result.MetaData)
	case err := <-errChan:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
	}
}

func (c *UserController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	resultChan := make(chan models.User)
	errChan := make(chan error)

	go func() {
		user, err := c.service.Get(ctx.Request().Context(), id)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- user
	}()

	select {
	case user := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusOK, user)
	case err := <-errChan:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user: "+err.Error())
	}
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

	return utils.SuccessResponse(ctx, http.StatusOK,  map[string]interface{}{"message": "User updated successfully"})
}

func (c *UserController) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

		if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete user")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, map[string]interface{}{"message": "User deleted successfully"})
}