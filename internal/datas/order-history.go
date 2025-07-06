package datas

import (
	"context"
	"zleeper-be/internal/models"

	"gorm.io/gorm"
)

type OrderHistoryData interface {
	Create(ctx context.Context, orderHistory *models.OrderHistory) error
	Update(ctx context.Context, orderHistory *models.OrderHistory) error
	CountData(ctx context.Context) (int, error)
	Fetch(ctx context.Context, limit int, offset int) ([]models.OrderHistory, error)
	GetByID(ctx context.Context, id uint) (models.OrderHistory, error)
	Delete(ctx context.Context, id uint) error
}

type orderHistoryData struct {
	db *gorm.DB
}

func (d *orderHistoryData) Create(ctx context.Context, orderHistory *models.OrderHistory) error {
	return d.db.WithContext(ctx).Create(orderHistory).Error
}

func (d *orderHistoryData) Update(ctx context.Context, orderHistory *models.OrderHistory) error {
	return d.db.WithContext(ctx).Save(orderHistory).Error
}

func (d *orderHistoryData) CountData(ctx context.Context) (int, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&models.OrderHistory{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return int(count), err
}

func (d *orderHistoryData) Fetch(ctx context.Context, limit int, offset int) ([]models.OrderHistory, error) {
	var orderItems []models.OrderHistory
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Preload("User").
		Preload("OrderItem").
		Limit(limit).
		Offset(offset).
		Find(&orderItems).Error
	return orderItems, err
}

func (d *orderHistoryData) GetByID(ctx context.Context, id uint) (models.OrderHistory, error) {
	var orderHistory models.OrderHistory
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Preload("User").
		Preload("OrderItem").
		First(&orderHistory, id).Error
	return orderHistory, err
}

func (d *orderHistoryData) Delete(ctx context.Context, id uint) error {
    return d.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&models.OrderHistory{}).
        Error
}