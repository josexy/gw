package model

import (
	"errors"
	"time"

	"github.com/josexy/gw/pkg/constants"
	"gorm.io/gorm"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info"`
	HTTPRule      *HttpRule      `json:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance"`
	AccessControl *AccessControl `json:"access_control"`
}

type ServiceInfo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"index;not null;size:130;column:service_name"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete"`
}

func (si *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (si *ServiceInfo) Find(db *gorm.DB) error {
	return db.Where(si).First(si).Error
}

func (si *ServiceInfo) Save(db *gorm.DB) error {
	return db.Save(si).Error
}

func (si *ServiceInfo) PageList(db *gorm.DB, info string, pageNo, pageSize int) ([]*ServiceInfo, int, error) {
	var list []*ServiceInfo
	var total int64
	offset := (pageNo - 1) * pageSize

	err := db.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&ServiceInfo{}).Where("is_delete = 0")
		if info != "" {
			query = query.Where("(service_name like ? or service_desc like ?)",
				"%"+info+"%", "%"+info+"%")
		}
		if err := query.Limit(pageSize).Offset(offset).Order("id desc").
			Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		return tx.Model(&ServiceInfo{}).Where("is_delete = 0").Count(&total).Error
	})
	return list, int(total), err
}

func (si *ServiceInfo) ServiceDetail(db *gorm.DB) (*ServiceDetail, error) {

	httpRule := &HttpRule{ServiceID: si.ID}
	tcpRule := &TcpRule{ServiceID: si.ID}
	grpcRule := &GrpcRule{ServiceID: si.ID}
	lbRule := &LoadBalance{ServiceID: si.ID}
	acRule := &AccessControl{ServiceID: si.ID}

	err := db.Transaction(func(tx *gorm.DB) error {
		var err error

		switch si.LoadType {
		case constants.LoadTypeHTTP:
			// 获取 HTTP 规则
			if err = httpRule.Find(tx); err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		case constants.LoadTypeTCP:
			// 获取 TCP 规则
			if err = tcpRule.Find(tx); err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		case constants.LoadTypeGRPC:
			// 获取 GRPC 规则
			if err = grpcRule.Find(tx); err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		default:
			return errors.New("load type not http/tcp/grpc")
		}

		// 获取 LoadBalance 规则
		if err = lbRule.Find(tx); err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// 获取 AccessControl 规则
		if err = acRule.Find(tx); err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ServiceDetail{
		Info:          si,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   lbRule,
		AccessControl: acRule,
	}, nil
}
