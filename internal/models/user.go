// internal/models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	FullName   string     `gorm:"not null" json:"full_name"`
	FirstOrder *time.Time `json:"first_order"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}