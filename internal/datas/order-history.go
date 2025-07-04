package datas

import (
	"gorm.io/gorm"
)

type OrderHistoryData interface {
	
}

type orderHistoryData struct {
	db *gorm.DB
}
