package services

import (
	"context"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/pkg/cache"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, page int, limit int) ([]models.User, error)
	Get(ctx context.Context, id uint) (models.User, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	data  datas.UserData
	cache cache.RedisCache
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	err := s.data.Create(ctx, user)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	
	err := s.data.Update(ctx, user)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_item:"+string(user.ID))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}

func (s *userService) List(ctx context.Context, page int, limit int) ([]models.User, error) {
	cacheKey := "order_items:page:" + string(page) + ":limit:" + string(limit)
	
	var cachedItems []models.User
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

func (s *userService) Get(ctx context.Context, id uint) (models.User, error) {
	cacheKey := "order_item:" + string(id)
	
	var cachedItem models.User
	if err := s.cache.Get(ctx, cacheKey, &cachedItem); err == nil {
		return cachedItem, nil
	}
	
	item, err := s.data.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	
	s.cache.Set(ctx, cacheKey, item, 5*time.Minute)
	
	return item, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	err := s.data.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "order_item:"+string(id))
	s.cache.Delete(ctx, "order_items:*")
	return nil
}