// internal/repositories/order_item_repository.go
package datas

import (
	"gorm.io/gorm"
)

type OrderItemData interface {
	
}

type orderItemData struct {
	db *gorm.DB
}

func NewOrderItemData(db *gorm.DB) OrderItemData {
	return &orderItemData{db: db}
}

