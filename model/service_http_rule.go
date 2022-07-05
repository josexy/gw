package model

import (
	"gorm.io/gorm"
)

type HttpRule struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	ServiceID      uint   `json:"service_id" gorm:"index;column:service_id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type"`
	Rule           string `json:"rule" gorm:"column:rule"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite"`
	HeaderTransfer string `json:"header_transfer" gorm:"column:header_transfer"`
}

func (rule *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (rule *HttpRule) Find(db *gorm.DB) error {
	return db.Where(rule).First(rule).Error
}

func (rule *HttpRule) Save(db *gorm.DB) error {
	return db.Save(rule).Error
}
