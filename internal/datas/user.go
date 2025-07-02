// internal/repositories/order_item_repository.go
package datas

import (
	"gorm.io/gorm"
)

type UserData interface {
	
}

type userData struct {
	db *gorm.DB
}

func NewUserData(db *gorm.DB) UserData {
	return &userData{db: db}
}