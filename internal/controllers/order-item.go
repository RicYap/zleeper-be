// internal/controllers/order_item_controller.go
package controllers

import (
	"zleeper-be/internal/services"
)

type OrderItemController struct {
	service services.OrderItemService
}

func NewOrderItemController(service services.OrderItemService) *OrderItemController {
	return &OrderItemController{service: service}
}

