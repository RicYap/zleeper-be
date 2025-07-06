package datas

import (
	"context"
	"zleeper-be/internal/models"

	"gorm.io/gorm"
)

type OrderItemData interface {
	Create(ctx context.Context, orderItem *models.OrderItem) error
	Update(ctx context.Context, orderItem *models.OrderItem) error
	CountData(ctx context.Context) (int, error)
	Fetch(ctx context.Context, limit int, offset int) ([]models.OrderItem, error)
	GetByID(ctx context.Context, id int) (models.OrderItem, error)
	Delete(ctx context.Context, id int) error
}

type orderItemData struct {
	db *gorm.DB
}

func (d *orderItemData) Create(ctx context.Context, orderItem *models.OrderItem) error {
	return d.db.WithContext(ctx).Create(orderItem).Error
}

func (d *orderItemData) Update(ctx context.Context, orderItem *models.OrderItem) error {
    return d.db.WithContext(ctx).
        Model(orderItem).
		Where("id = ?", orderItem.ID).
        UpdateColumns(map[string]interface{}{
            "name":  orderItem.Name,
            "price":    orderItem.Price,
            "expired_at":  orderItem.ExpiredAt,
            "updated_at":  orderItem.UpdatedAt, 
        }).Error
}

func (d *orderItemData) CountData(ctx context.Context) (int, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&models.OrderItem{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return int(count), err
}

func (d *orderItemData) Fetch(ctx context.Context, limit int, offset int) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&orderItems).Error
	return orderItems, err
}

func (d *orderItemData) GetByID(ctx context.Context, id int) (models.OrderItem, error) {
	var orderItem models.OrderItem
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&orderItem, id).Error
	return orderItem, err
}

func (d *orderItemData) Delete(ctx context.Context, id int) error {
    return d.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&models.OrderItem{}).
        Error
}