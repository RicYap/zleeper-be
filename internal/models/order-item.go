// internal/models/order_item.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"not null" json:"name"`
    Price     float64        `gorm:"not null" json:"price"`
    ExpiredAt time.Time      `gorm:"not null" json:"expired_at"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type OrderItemPagination struct {
	Data 		[]OrderItem 	`json:"data"`
	MetaData 	any	 			`json:"metadata"`
}