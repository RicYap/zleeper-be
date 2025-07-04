// internal/services/order_item_service.go
package services

import (
	"zleeper-be/internal/datas"
	"zleeper-be/pkg/cache"
)

type UserService interface {
	
}

type userService struct {
	repo  datas.UserData
	cache cache.RedisCache
}

