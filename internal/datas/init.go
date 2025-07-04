package datas

import "gorm.io/gorm"

type Datas struct {
    OrderItem    OrderItemData
    OrderHistory OrderHistoryData
    User         UserData
}

func NewDatas(db *gorm.DB) *Datas {
    return &Datas{
        OrderItem:    &orderItemData{db: db},
        OrderHistory: &orderHistoryData{db: db},
        User:         &userData{db: db},
    }
}