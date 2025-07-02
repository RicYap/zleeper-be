// internal/services/order_item_service.go
package services

import (
	"zleeper-be/internal/datas"
	"zleeper-be/pkg/cache"
)

type OrderItemService interface {
	
}

type orderItemService struct {
	repo  datas.OrderItemData
	cache cache.RedisCache
}

func NewOrderItemService(repo datas.OrderItemData, cache *cache.RedisCache) OrderItemService {
	return &orderItemService{repo: repo, cache: *cache}
}

