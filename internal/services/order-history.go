// internal/services/order_item_service.go
package services

import (
	"context"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/pkg/cache"
)

type OrderHistoryService interface {
	Create(ctx context.Context, orderHistory *models.OrderHistory) error
	Update(ctx context.Context, orderHistory *models.OrderHistory) error
	List(ctx context.Context, page int, limit int) ([]models.OrderHistory, error)
	Get(ctx context.Context, id uint) (models.OrderHistory, error)
	Delete(ctx context.Context, id uint) error
}

type orderHistoryService struct {
	data  datas.OrderHistoryData
	userService UserService
	cache cache.RedisCache
}

func (s *orderHistoryService) Create(ctx context.Context, orderHistory *models.OrderHistory) error {
	orderHistory.CreatedAt = time.Now()
	orderHistory.UpdatedAt = time.Now()

	err := s.userService.MarkFirstOrder(ctx, orderHistory.UserID, orderHistory.CreatedAt)
	if err != nil {
		return err
	}

	err = s.data.Create(ctx, orderHistory)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_histories:*")
	return nil
}

func (s *orderHistoryService) Update(ctx context.Context, orderHistory *models.OrderHistory) error {
	orderHistory.UpdatedAt = time.Now()
	
	err := s.data.Update(ctx, orderHistory)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_history:"+string(orderHistory.ID))
	s.cache.Delete(ctx, "order_histories:*")
	return nil
}

func (s *orderHistoryService) List(ctx context.Context, page int, limit int) ([]models.OrderHistory, error) {
	cacheKey := "order_histories:page:" + string(page) + ":limit:" + string(limit)
	
	var cachedItems []models.OrderHistory
	if err := s.cache.Get(ctx, cacheKey, &cachedItems); err == nil {
		return cachedItems, nil
	}
	
	offset := (page - 1) * limit
	
	items, err := s.data.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	s.cache.Set(ctx, cacheKey, items, 5*time.Minute)
	
	return items, nil
}

func (s *orderHistoryService) Get(ctx context.Context, id uint) (models.OrderHistory, error) {
	cacheKey := "order_history:" + string(id)
	
	var cachedItem models.OrderHistory
	if err := s.cache.Get(ctx, cacheKey, &cachedItem); err == nil {
		return cachedItem, nil
	}
	
	item, err := s.data.GetByID(ctx, id)
	if err != nil {
		return models.OrderHistory{}, err
	}
	
	s.cache.Set(ctx, cacheKey, item, 5*time.Minute)
	
	return item, nil
}

func (s *orderHistoryService) Delete(ctx context.Context, id uint) error {
	err := s.data.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_history:"+string(id))
	s.cache.Delete(ctx, "order_histories:*")
	return nil
}

// func (s *orderHistoryService) toResponse(orderHistory *models.OrderHistory) *models.OrderHistoryResponse {
// 	if orderHistory == nil {
// 		return nil
// 	}

// 	return &models.OrderHistoryResponse{
// 		ID:           orderHistory.ID,
// 		UserID:       orderHistory.UserID,
// 		OrderItemID:  orderHistory.OrderItemID,
// 		Descriptions: orderHistory.Description,
// 		CreatedAt:    orderHistory.CreatedAt,
// 		UserName:     orderHistory.User.FullName,
// 		ItemName:     orderHistory.OrderItem.Name,
// 		ItemPrice:    orderHistory.OrderItem.Price,
// 	}
// }