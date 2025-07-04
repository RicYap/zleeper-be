// internal/services/order_item_service.go
package services

import (
	"zleeper-be/internal/datas"
	"zleeper-be/pkg/cache"
)

type OrderHistoryService interface {
	
}

type orderHistoryService struct {
	repo  datas.OrderHistoryData
	cache cache.RedisCache
}

