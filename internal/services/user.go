package services

import (
	"context"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/internal/utils"
	"zleeper-be/pkg/cache"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, page int, limit int) (models.UserPagination, error)
	Get(ctx context.Context, id uint) (models.User, error)
	Delete(ctx context.Context, id uint) error
	MarkFirstOrder(ctx context.Context, userID uint, orderTime time.Time) error
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
	
	s.cache.Delete(ctx, "users:*")
	return nil
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	
	err := s.data.Update(ctx, user)
	if err != nil {
		return err
	}
	
	s.cache.Delete(ctx, "user:"+string(rune(user.ID)))
	s.cache.Delete(ctx, "users:*")
	return nil
}

func (s *userService) List(ctx context.Context, page int, limit int) (models.UserPagination, error) {
	cacheKey := "users:page:" + string(rune(page)) + ":limit:" + string(rune(limit))
	
	var cachedItems models.UserPagination
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

func (s *userService) Get(ctx context.Context, id uint) (models.User, error) {
	cacheKey := "user:" + string(rune(id))
	
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
	
	s.cache.Delete(ctx, "user:"+string(rune(id)))
	s.cache.Delete(ctx, "users:*")
	return nil
}


func (s *userService) MarkFirstOrder(ctx context.Context, userID uint, orderTime time.Time) error {
	
	user, err := s.Get(ctx, userID)
	if err != nil {
		return err
	}

	if user.FirstOrder != nil {
		return nil
	}

	user.FirstOrder = &orderTime

	return s.Update(ctx, &user)
}