// internal/models/order_history.go
package models

import "time"

type OrderHistory struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	OrderItemID uint       `gorm:"not null" json:"order_item_id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        User       `gorm:"foreignKey:UserID" json:"user"`
	OrderItem   OrderItem  `gorm:"foreignKey:OrderItemID" json:"order_item"`
}