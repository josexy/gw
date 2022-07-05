package model

import (
	"strings"

	"gorm.io/gorm"
)

type LoadBalance struct {
	ID                     uint   `json:"id" gorm:"primaryKey"`
	ServiceID              uint   `json:"service_id" gorm:"index;column:service_id"`
	RoundType              int    `json:"round_type" gorm:"column:round_type"`
	IpList                 string `json:"ip_list" gorm:"column:ip_list"`
	WeightList             string `json:"weight_list" gorm:"column:weight_list"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" gorm:"column:upstream_max_idle"`
}

func (lb *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (lb *LoadBalance) Find(db *gorm.DB) error {
	return db.Where(lb).First(lb).Error
}

func (lb *LoadBalance) Save(db *gorm.DB) error {
	return db.Save(lb).Error
}

func (lb *LoadBalance) GetIPList() []string {
	return strings.Split(lb.IpList, ",")
}

func (lb *LoadBalance) GetWeightList() []string {
	return strings.Split(lb.WeightList, ",")
}
