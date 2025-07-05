// internal/models/order_history.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderHistory struct {
	ID          uint       		`gorm:"primaryKey" json:"id"`
	UserID      uint       		`gorm:"not null" json:"user_id"`
	OrderItemID uint       		`gorm:"not null" json:"order_item_id"`
	Description string     		`json:"description"`
	CreatedAt   time.Time  		`json:"created_at"`
	UpdatedAt   time.Time  		`json:"updated_at"`
	DeletedAt 	gorm.DeletedAt 	`gorm:"index" json:"deleted_at,omitempty"`
	User        User       		`gorm:"foreignKey:UserID" json:"user"`
	OrderItem   OrderItem  		`gorm:"foreignKey:OrderItemID" json:"order_item"`
}

type OrderHistoryResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	OrderItemID  uint      `json:"order_item_id"`
	Descriptions string    `json:"descriptions"`
	CreatedAt    time.Time `json:"created_at"`
	UserName     string    `json:"user_name"`
	ItemName     string    `json:"item_name"`
	ItemPrice    float64   `json:"item_price"`
}