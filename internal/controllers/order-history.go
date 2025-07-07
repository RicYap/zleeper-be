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

type OrderHistoryController struct {
	service services.OrderHistoryService
}

func (c *OrderHistoryController) Create(ctx echo.Context) error {
	var orderHistory models.OrderHistory
	if err := ctx.Bind(&orderHistory); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
	}

	resultChan := make(chan uint)
	errChan := make(chan error)

	go func() {
		if err := c.service.Create(ctx.Request().Context(), &orderHistory); err != nil {
			errChan <- err
			return
		}
		resultChan <- orderHistory.ID
	}()

	select {
	case id := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusCreated, map[string]interface{}{
			"id": id,
		})
	case err := <-errChan:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create order history: "+err.Error())
	}
}

func (c *OrderHistoryController) List(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	resultChan := make(chan models.OrderHistoryPagination)
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
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch order histories: "+err.Error())
	}
}

func (c *OrderHistoryController) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	resultChan := make(chan models.OrderHistoryResponse)
	errChan := make(chan error)

	go func() {
		order, err := c.service.Get(ctx.Request().Context(), id)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- order
	}()

	select {
	case order := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusOK, order)
	case err := <-errChan:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Order history not found")
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get order history: "+err.Error())
	}
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

	return utils.SuccessResponse(ctx, http.StatusOK, map[string]interface{}{"message": "Order history updated successfully"})
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

func (c *OrderHistoryController) GetAllOrdersByDate(ctx echo.Context) error {
	resultChan := make(chan []models.OrdersByDateResponse)
	errChan := make(chan error)

	go func() {
		orders, err := c.service.GetAllOrdersGroupedByDate(ctx.Request().Context())
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- orders
	}()

	select {
	case orders := <-resultChan:
		return utils.SuccessResponse(ctx, http.StatusOK, orders)
	case err := <-errChan:
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get order history by date: "+err.Error())
	}
}