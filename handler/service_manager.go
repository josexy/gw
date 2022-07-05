package handler

import (
	"errors"
	"strings"
	"sync"

	"github.com/josexy/gw/global"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
)

var ServiceManagerHandler *ServiceManager

type ServiceManager struct {
	handler MapBaseHandler[*model.ServiceDetail]
	init    sync.Once
	err     error
}

func init() {
	ServiceManagerHandler = NewServiceManager()
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		handler: NewMapBaseHandler[*model.ServiceDetail](),
	}
}

func (sm *ServiceManager) GetTcpServiceList() []*model.ServiceDetail {
	sm.handler.RLock()
	defer sm.handler.RUnlock()
	var list []*model.ServiceDetail
	for _, serverItem := range sm.handler.List {
		if serverItem.Info.LoadType == constants.LoadTypeTCP {
			list = append(list, serverItem)
		}
	}
	return list
}

func (sm *ServiceManager) GetGrpcServiceList() []*model.ServiceDetail {
	sm.handler.RLock()
	defer sm.handler.RUnlock()
	var list []*model.ServiceDetail
	for _, serverItem := range sm.handler.List {
		if serverItem.Info.LoadType == constants.LoadTypeGRPC {
			list = append(list, serverItem)
		}
	}
	return list
}

// HTTPAccessMode 请求的URL/HOST是否匹配Service规则
func (sm *ServiceManager) HTTPAccessMode(host, path string) (*model.ServiceDetail, error) {
	//1、前缀匹配 /abc
	//2、域名匹配 www.test.com

	// www.test.com:8080/abc/def?id=1&age=12
	// host: 127.0.0.1:8080/www.test.com:8080
	// url: /abc
	host = host[0:strings.Index(host, ":")]

	sm.handler.RLock()
	defer sm.handler.RUnlock()

	for _, item := range sm.handler.List {
		if item.Info.LoadType != constants.LoadTypeHTTP {
			continue
		}

		// 是否匹配服务规则
		switch item.HTTPRule.RuleType {
		case constants.HTTPRuleTypeDomain: // 域名匹配
			if item.HTTPRule.Rule == host {
				return item, nil
			}
		case constants.HTTPRuleTypePrefixURL: // 前缀匹配
			if strings.HasPrefix(path, item.HTTPRule.Rule) {
				return item, nil
			}
		}
	}
	return nil, errors.New("no matched service")
}

func (sm *ServiceManager) LoadOnce() error {
	sm.init.Do(func() {
		var serviceInfo model.ServiceInfo
		// 获取所有服务信息
		list, _, err := serviceInfo.PageList(global.DB, "", 1, 99999)
		if err != nil {
			sm.err = err
			return
		}

		sm.handler.Lock()
		defer sm.handler.Unlock()

		for _, item := range list {
			serviceDetail, err := item.ServiceDetail(global.DB)
			if err != nil {
				sm.err = err
				return
			}
			sm.handler.Cache[item.ServiceName] = serviceDetail
			sm.handler.List = append(sm.handler.List, serviceDetail)
		}
	})
	return sm.err
}
