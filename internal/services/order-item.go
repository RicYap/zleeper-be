// internal/services/order_item_service.go
package services

import (
	"context"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/pkg/cache"
)

type OrderItemService interface {
	Create(ctx context.Context, orderItem *models.OrderItem) error
	Update(ctx context.Context, orderItem *models.OrderItem) error
	List(ctx context.Context, page int, limit int) ([]models.OrderItem, error)
	Get(ctx context.Context, id uint) (models.OrderItem, error)
	Delete(ctx context.Context, id uint) error
}

type orderItemService struct {
	data  datas.OrderItemData
	cache cache.RedisCache
}

func (s *orderItemService) Create(ctx context.Context, orderItem *models.OrderItem) error {
	orderItem.CreatedAt = time.Now()
	orderItem.UpdatedAt = time.Now()
	
	err := s.data.Create(ctx, orderItem)
	if err != nil {
		return err
	}
	
	// Invalidate cache
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *orderItemService) Update(ctx context.Context, orderItem *models.OrderItem) error {
	orderItem.UpdatedAt = time.Now()
	
	err := s.data.Update(ctx, orderItem)
	if err != nil {
		return err
	}
	
	// Invalidate cache for this item and list
	s.cache.Delete(ctx, "order_item:"+string(orderItem.ID))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *orderItemService) List(ctx context.Context, page int, limit int) ([]models.OrderItem, error) {
	cacheKey := "order_items:page:" + string(page) + ":limit:" + string(limit)
	
	// Try to get from cache first
	var cachedItems []models.OrderItem
	if err := s.cache.Get(ctx, cacheKey, &cachedItems); err == nil {
		return cachedItems, nil
	}
	
	// Calculate offset
	offset := (page - 1) * limit
	
	// Get from database
	items, err := s.data.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	s.cache.Set(ctx, cacheKey, items, 5*time.Minute)
	
	return items, nil
}

func (s *orderItemService) Get(ctx context.Context, id uint) (models.OrderItem, error) {
	cacheKey := "order_item:" + string(id)
	
	// Try to get from cache first
	var cachedItem models.OrderItem
	if err := s.cache.Get(ctx, cacheKey, &cachedItem); err == nil {
		return cachedItem, nil
	}
	
	// Get from database
	item, err := s.data.GetByID(ctx, id)
	if err != nil {
		return models.OrderItem{}, err
	}
	
	// Cache the result
	s.cache.Set(ctx, cacheKey, item, 5*time.Minute)
	
	return item, nil
}

func (s *orderItemService) Delete(ctx context.Context, id uint) error {
	err := s.data.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	// Invalidate cache for this item and list
	s.cache.Delete(ctx, "order_item:"+string(id))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}