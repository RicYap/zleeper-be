package datas

import (
	"context"
	"zleeper-be/internal/models"

	"gorm.io/gorm"
)

type UserData interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	CountData(ctx context.Context) (int, error)
	Fetch(ctx context.Context, limit int, offset int) ([]models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	Delete(ctx context.Context, id int) error
}

type userData struct {
	db *gorm.DB
}

func (d *userData) Create(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userData) Update(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).
        Model(user).
		Where("id = ?", user.ID).
        UpdateColumns(map[string]interface{}{
            "full_name":  user.FullName,
            "updated_at":  user.UpdatedAt, 
        }).Error
}

func (d *userData) CountData(ctx context.Context) (int, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return int(count), err
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

func (d *userData) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		First(&user, id).Error
	return user, err
}

func (d *userData) Delete(ctx context.Context, id int) error {
    return d.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&models.User{}).
        Error
}