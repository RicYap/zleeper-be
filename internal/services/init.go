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
	return &Services{
		OrderItem:    NewOrderItemService(datas.OrderItem, redisCache),
		OrderHistory: NewOrderHistoryService(datas.OrderHistory, redisCache),
		User:         NewUserService(datas.User, redisCache),
	}
}

func NewOrderItemService(data datas.OrderItemData, cache *cache.RedisCache) OrderItemService {
	return &orderItemService{data: data, cache: *cache}
}

func NewOrderHistoryService(data datas.OrderHistoryData, cache *cache.RedisCache) OrderHistoryService {
	return &orderHistoryService{data: data, cache: *cache}
}

func NewUserService(data datas.UserData, cache *cache.RedisCache) UserService {
	return &userService{data: data, cache: *cache}
}