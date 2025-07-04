package datas

import (
	"gorm.io/gorm"
)

type UserData interface {
	
}

type userData struct {
	db *gorm.DB
}
