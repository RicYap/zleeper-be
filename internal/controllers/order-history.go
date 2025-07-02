package controllers

import (
	"zleeper-be/internal/services"
)

type OrderHistoryController struct {
	service services.OrderHistoryService
}

func NewOrderHistoryController(service services.OrderHistoryService) *OrderHistoryController {
	return &OrderHistoryController{service: service}
}
