package model

import (
	"gorm.io/gorm"
)

type GrpcRule struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	ServiceID      uint   `json:"service_id" gorm:"index;column:service_id"`
	Port           int    `json:"port" gorm:"column:port"`
	HeaderTransfer string `json:"header_transfer" gorm:"column:header_transfer"`
}

func (rule *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (rule *GrpcRule) Find(db *gorm.DB) error {
	return db.Where(rule).First(rule).Error
}

func (rule *GrpcRule) Save(db *gorm.DB) error {
	return db.Save(rule).Error
}
