package services

import (
	"zleeper-be/internal/datas"
	"zleeper-be/pkg/cache"
)

type Services struct {
	OrderItem    OrderItemService
	OrderHistory OrderHistoryService
	User         UserService
}

func NewServices(datas *datas.Datas, redisCache *cache.RedisCache) *Services {

	orderItemService :=  NewOrderItemService(datas.OrderItem, redisCache)
	userService :=  NewUserService(datas.User, redisCache)
	orderHistoryService := NewOrderHistoryService(datas.OrderHistory, userService, redisCache)
	
	return &Services{
		OrderItem:    orderItemService,
		OrderHistory: orderHistoryService,
		User:         userService,
	}
}

func NewOrderItemService(data datas.OrderItemData, cache *cache.RedisCache) OrderItemService {
	return &orderItemService{data: data, cache: *cache}
}

func NewOrderHistoryService(data datas.OrderHistoryData, userService UserService, cache *cache.RedisCache) OrderHistoryService {
	return &orderHistoryService{data: data, userService: userService , cache: *cache}
}

func NewUserService(data datas.UserData, cache *cache.RedisCache) UserService {
	return &userService{data: data, cache: *cache}
}