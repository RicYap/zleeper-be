package datas

import (
	"context"
	"zleeper-be/internal/models"

	"gorm.io/gorm"
)

type UserData interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Fetch(ctx context.Context, limit int, offset int) ([]models.User, error)
	GetByID(ctx context.Context, id uint) (models.User, error)
	Delete(ctx context.Context, id uint) error
}

type userData struct {
	db *gorm.DB
}

func (d *userData) Create(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userData) Update(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Save(user).Error
}

func (d *userData) Fetch(ctx context.Context, limit int, offset int) ([]models.User, error) {
	var orderItems []models.User
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&orderItems).Error
	return orderItems, err
}

func (d *userData) GetByID(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&user, id).Error
	return user, err
}

func (d *userData) Delete(ctx context.Context, id uint) error {
    return d.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&models.User{}).
        Error
}