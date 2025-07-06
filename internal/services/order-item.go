// internal/services/order_item_service.go
package services

import (
	"context"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/internal/utils"
	"zleeper-be/pkg/cache"
)

type OrderItemService interface {
	Create(ctx context.Context, orderItem *models.OrderItem) error
	Update(ctx context.Context, orderItem *models.OrderItem) error
	List(ctx context.Context, page int, limit int) (models.OrderItemPagination, error)
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
	
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *orderItemService) Update(ctx context.Context, orderItem *models.OrderItem) error {
	orderItem.UpdatedAt = time.Now()
	
	err := s.data.Update(ctx, orderItem)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_item:"+string(rune(orderItem.ID)))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *orderItemService) List(ctx context.Context, page int, limit int) (models.OrderItemPagination, error) {
	cacheKey := "order_items:page:" + string(rune(page)) + ":limit:" + string(rune(limit))
	
	var cachedItems models.OrderItemPagination
	if err := s.cache.Get(ctx, cacheKey, &cachedItems); err == nil {
		return cachedItems, nil
	}
	
	totalData, err := s.data.CountData(ctx)
	if err != nil {
		return cachedItems, err
	}

	offset, metaData := utils.PaginationCalculation(page, limit, totalData)
	
	items, err := s.data.Fetch(ctx, limit, offset)
	if err != nil {
		return cachedItems, err
	}

	cachedItems.Data = items
	cachedItems.MetaData = metaData
	
	s.cache.Set(ctx, cacheKey, items, 5*time.Minute)
	
	return cachedItems, nil
}

func (s *orderItemService) Get(ctx context.Context, id uint) (models.OrderItem, error) {
	cacheKey := "order_item:" + string(rune(id))
	
	var cachedItem models.OrderItem
	if err := s.cache.Get(ctx, cacheKey, &cachedItem); err == nil {
		return cachedItem, nil
	}
	
	item, err := s.data.GetByID(ctx, id)
	if err != nil {
		return models.OrderItem{}, err
	}
	
	s.cache.Set(ctx, cacheKey, item, 5*time.Minute)
	
	return item, nil
}

func (s *orderItemService) Delete(ctx context.Context, id uint) error {
	err := s.data.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_item:"+string(rune(id)))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}