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

type OrderItemController struct {
	service services.OrderItemService
}

func (c *OrderItemController) Create(ctx echo.Context) error {
	var orderItem models.OrderItem
	if err := ctx.Bind(&orderItem); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	resultChan := make(chan uint)
	errChan := make(chan error)

	go func() {
		if err := c.service.Create(ctx.Request().Context(), &orderItem); err != nil {
			errChan <- err
			return
		}
		resultChan <- orderItem.ID
	}()

	select {
	case id := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusCreated, map[string]interface{}{
			"id": id,
		})
	case err := <-errChan:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create order item: "+err.Error())
	}
}

func (c *OrderItemController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resultChan := make(chan models.OrderItemPagination)
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
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch order items: "+err.Error())
	}
}

func (c *OrderItemController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	resultChan := make(chan models.OrderItem)
	errChan := make(chan error)

	go func() {
		item, err := c.service.Get(ctx.Request().Context(), id)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- item
	}()

	select {
	case item := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusOK, item)
	case err := <-errChan:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Order item not found")
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get order item: "+err.Error())
	}
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