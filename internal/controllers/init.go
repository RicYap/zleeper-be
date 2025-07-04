package controllers

import "zleeper-be/internal/services"

type Controllers struct {
	OrderItem    *OrderItemController
	OrderHistory *OrderHistoryController
	User         *UserController
}

func NewControllers(s *services.Services) *Controllers {
	return &Controllers{
		OrderItem:    NewOrderItemController(s.OrderItem),
		OrderHistory: NewOrderHistoryController(s.OrderHistory),
		User:         NewUserController(s.User),
	}
}

func NewOrderHistoryController(service services.OrderHistoryService) *OrderHistoryController {
	return &OrderHistoryController{service: service}
}

func NewOrderItemController(service services.OrderItemService) *OrderItemController {
	return &OrderItemController{service: service}
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}
