package model

import (
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	ServiceID uint `json:"service_id" gorm:"index;column:service_id"`
	Port      int  `json:"port" gorm:"column:port"`
}

func (rule *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (rule *TcpRule) Find(db *gorm.DB) error {
	return db.Where(rule).First(rule).Error
}

func (rule *TcpRule) Save(db *gorm.DB) error {
	return db.Save(rule).Error
}
