// internal/controllers/order_item_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"zleeper-be/internal/models"
	"zleeper-be/internal/services"
	"zleeper-be/internal/utils"

	"github.com/labstack/echo/v4"
)

type OrderItemController struct {
	service services.OrderItemService
}

func (c *OrderItemController) Create(ctx echo.Context) error {
	var orderItem models.OrderItem
	if err := ctx.Bind(&orderItem); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.service.Create(ctx.Request().Context(), &orderItem); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create order item")
	}

	return utils.SuccessResponse(ctx, http.StatusCreated,  map[string]interface{}{"message": "Order item created successfully"})
}

func (c *OrderItemController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	orderItemsPagination, err := c.service.List(ctx.Request().Context(), page, limit)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch order items")
	}

	return utils.PaginatedResponse(ctx, http.StatusOK, orderItemsPagination.Data, orderItemsPagination.MetaData)
}

func (c *OrderItemController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	orderItem, err := c.service.Get(ctx.Request().Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Order item not found")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, orderItem)
}

func (c *OrderItemController) Update(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	var orderItem models.OrderItem
	if err := ctx.Bind(&orderItem); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	orderItem.ID = uint(id)
	if err := c.service.Update(ctx.Request().Context(), &orderItem); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update order item")
	}

	return utils.SuccessResponse(ctx, http.StatusOK,  map[string]interface{}{"message": "Order item updated successfully"})
}

func (c *OrderItemController) Delete(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete order item")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, map[string]interface{}{"message": "Order item deleted successfully"})
}