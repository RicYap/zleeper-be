package controllers

import (
	"net/http"
	"strconv"
	"zleeper-be/internal/models"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"

	"github.com/labstack/echo/v4"
)

type OrderHistoryController struct {
	service services.OrderHistoryService
}

func (c *OrderHistoryController) Create(ctx echo.Context) error {
	var orderHistory models.OrderHistory
	if err := ctx.Bind(&orderHistory); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.service.Create(ctx.Request().Context(), &orderHistory); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create order history")
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, map[string]interface{}{"message": "Order history created successfully"})
}

func (c *OrderHistoryController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	orderHistoriesPagination, err := c.service.List(ctx.Request().Context(), page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch order histories")
	}

	return utils.PaginatedResponse(ctx, http.StatusOK, orderHistoriesPagination.Data, orderHistoriesPagination.MetaData)
}

func (c *OrderHistoryController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	orderHistory, err := c.service.Get(ctx.Request().Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Order history not found")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, orderHistory)
}

func (c *OrderHistoryController) Update(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	var orderHistory models.OrderHistory
	if err := ctx.Bind(&orderHistory); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	orderHistory.ID = uint(id)
	if err := c.service.Update(ctx.Request().Context(), &orderHistory); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update order history")
	}

	return utils.SuccessResponse(ctx, http.StatusOK,  map[string]interface{}{"message": "Order history updated successfully"})
}

func (c *OrderHistoryController) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete order history")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, map[string]interface{}{"message": "Order history deleted successfully"})
}
