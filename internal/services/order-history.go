// internal/services/order_item_service.go
package services

import (
	"context"
	"sort"
	"strconv"
	"time"
	"zleeper-be/internal/datas"
	"zleeper-be/internal/models"
	"zleeper-be/internal/utils"
	"zleeper-be/pkg/cache"
)

type OrderHistoryService interface {
	Create(ctx context.Context, orderHistory *models.OrderHistory) error
	Update(ctx context.Context, orderHistory *models.OrderHistory) error
	List(ctx context.Context, page int, limit int) (models.OrderHistoryPagination, error)
	Get(ctx context.Context, id int) (models.OrderHistoryResponse, error)
	Delete(ctx context.Context, id int) error
	GetAllOrdersGroupedByDate(ctx context.Context) ([]models.OrdersByDateResponse, error)
}

type orderHistoryService struct {
	data        datas.OrderHistoryData
	userService UserService
	cache       cache.RedisCache
}

func (s *orderHistoryService) Create(ctx context.Context, orderHistory *models.OrderHistory) error {
	orderHistory.CreatedAt = time.Now()
	orderHistory.UpdatedAt = time.Now()

	err := s.userService.MarkFirstOrder(ctx, int(orderHistory.UserID), orderHistory.CreatedAt)
	if err != nil {
		return err
	}

	err = s.data.Create(ctx, orderHistory)
	if err != nil {
		return err
	}

	s.cache.DeleteAll(ctx, "order_histories:*")
	return nil
}

func (s *orderHistoryService) Update(ctx context.Context, orderHistory *models.OrderHistory) error {

	idString := strconv.Itoa(int(orderHistory.ID))

	orderHistory.UpdatedAt = time.Now()

	err := s.data.Update(ctx, orderHistory)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, "order_history:"+idString)
	s.cache.DeleteAll(ctx, "order_histories:*")
	return nil
}

func (s *orderHistoryService) List(ctx context.Context, page int, limit int) (models.OrderHistoryPagination, error) {

	var (
		cachedItems   models.OrderHistoryPagination
		itemResponses []models.OrderHistoryResponse
	)

	pageString := strconv.Itoa(page)
	limitString := strconv.Itoa(limit)

	cacheKey := "order_histories:page:" + pageString + ":limit:" + limitString

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

	for _, item := range items {
		itemResponse := s.FormatResponse(item)

		itemResponses = append(itemResponses, itemResponse)
	}

	cachedItems.Data = itemResponses
	cachedItems.MetaData = metaData

	s.cache.Set(ctx, cacheKey, cachedItems, 5*time.Minute)

	return cachedItems, nil
}

func (s *orderHistoryService) Get(ctx context.Context, id int) (models.OrderHistoryResponse, error) {

	idString := strconv.Itoa(id)

	cacheKey := "order_history:" + idString

	var cachedItem models.OrderHistoryResponse
	if err := s.cache.Get(ctx, cacheKey, &cachedItem); err == nil {
		return cachedItem, nil
	}

	item, err := s.data.GetByID(ctx, id)
	if err != nil {
		return models.OrderHistoryResponse{}, err
	}

	itemResponse := s.FormatResponse(item)

	s.cache.Set(ctx, cacheKey, itemResponse, 5*time.Minute)

	return itemResponse, nil
}

func (s *orderHistoryService) Delete(ctx context.Context, id int) error {

	idString := strconv.Itoa(id)

	err := s.data.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, "order_history:"+idString)
	s.cache.DeleteAll(ctx, "order_histories:*")
	return nil
}

func (s *orderHistoryService) GetAllOrdersGroupedByDate(ctx context.Context) ([]models.OrdersByDateResponse, error) {

	cacheKey := "order_histories:by_date:all"
	var cachedResponse []models.OrdersByDateResponse
	if err := s.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		return cachedResponse, nil
	}

	ordersByDate, err := s.data.GetAllGroupedByDate(ctx)
	if err != nil {
		return nil, err
	}

	var response []models.OrdersByDateResponse
	for date, orders := range ordersByDate {
		response = append(response, models.OrdersByDateResponse{
			Date:   date.Format("2006-01-02"),
			Orders: s.FormatResponseList(orders),
		})
	}

	sort.Slice(response, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", response[i].Date)
		dateJ, _ := time.Parse("2006-01-02", response[j].Date)
		return dateI.After(dateJ)
	})

	s.cache.Set(ctx, cacheKey, response, 10*time.Minute)

	return response, nil
}

func (s *orderHistoryService) FormatResponseList(histories []models.OrderHistory) []models.OrderHistoryResponse {
	responses := make([]models.OrderHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.FormatResponse(history)
	}
	return responses
}

func (s *orderHistoryService) FormatResponse(orderHistory models.OrderHistory) models.OrderHistoryResponse {

	return models.OrderHistoryResponse{
		ID:           orderHistory.ID,
		UserID:       orderHistory.UserID,
		OrderItemID:  orderHistory.OrderItemID,
		Descriptions: orderHistory.Description,
		CreatedAt:    orderHistory.CreatedAt,
		UserName:     orderHistory.User.FullName,
		ItemName:     orderHistory.OrderItem.Name,
		ItemPrice:    orderHistory.OrderItem.Price,
	}
}
