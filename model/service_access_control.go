package model

import (
	"gorm.io/gorm"
)

type AccessControl struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	ServiceID         uint   `json:"service_id" gorm:"index;column:service_id"`
	EnableAuth        int    `json:"enable_auth" gorm:"column:enable_auth"`
	BlackList         string `json:"black_list" gorm:"column:black_list"`
	WhiteList         string `json:"white_list" gorm:"column:white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit"`
}

func (ac *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (ac *AccessControl) Find(db *gorm.DB) error {
	return db.Where(ac).First(ac).Error
}

func (ac *AccessControl) Save(db *gorm.DB) error {
	return db.Save(ac).Error
}
