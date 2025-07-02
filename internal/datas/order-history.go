// internal/repositories/order_item_repository.go
package datas

import (
	"gorm.io/gorm"
)

type OrderHistoryData interface {
	
}

type orderHistoryData struct {
	db *gorm.DB
}

func NewOrderHistoryData(db *gorm.DB) OrderHistoryData {
	return &orderHistoryData{db: db}
}