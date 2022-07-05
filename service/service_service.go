package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josexy/gw/global"
	"github.com/josexy/gw/handler"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/codes"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/serializer"
	"gorm.io/gorm"
)

type ListService struct {
	Info     string `json:"info" form:"info"`
	PageNo   int    `json:"page_no" form:"page_no" binding:"required,min=0"`
	PageSize int    `json:"page_size" form:"page_size" binding:"required,min=0"`
}

type DetailService struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type DeleteService struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type StatService struct {
	ID int `json:"id" form:"id" binding:"required"`
}

type AddHTTPService struct {
	ServiceName string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc string `json:"service_desc" binding:"required,max=255,min=1"`

	RuleType       int    `json:"rule_type" binding:"max=1,min=0"`
	Rule           string `json:"rule" binding:"required,valid_rule"`
	NeedHttps      int    `json:"need_https" binding:"max=1,min=0"`
	NeedStripUri   int    `json:"need_strip_uri" binding:"max=1,min=0"`
	UrlRewrite     string `json:"url_rewrite" binding:"valid_url_rewrite"`
	HeaderTransfer string `json:"header_transfer" binding:"valid_header_transfer"`

	EnableAuth        int    `json:"enable_auth" validate:"max=1,min=0"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" binding:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" binding:"min=0"`

	RoundType              int    `json:"round_type" binding:"max=3,min=0"`
	IpList                 string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList             string `json:"weight_list" binding:"required,valid_weightlist"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" binding:"min=0"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" binding:"min=0"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" binding:"min=0"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" binding:"min=0"`
}

type UpdateHTTPService struct {
	ID          int    `json:"id" binding:"required,min=1"`
	ServiceName string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc string `json:"service_desc" binding:"required,max=255,min=1"`

	RuleType       int    `json:"rule_type" binding:"max=1,min=0"`
	Rule           string `json:"rule" binding:"required,valid_rule"`
	NeedHttps      int    `json:"need_https" binding:"max=1,min=0"`
	NeedStripUri   int    `json:"need_strip_uri" binding:"max=1,min=0"`
	UrlRewrite     string `json:"url_rewrite" binding:"valid_url_rewrite"`
	HeaderTransfer string `json:"header_transfer" binding:"valid_header_transfer"`

	EnableAuth        int    `json:"enable_auth" validate:"max=1,min=0"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" binding:"min=0"`
	ServiceFlowLimit  int    `json:"service_flow_limit" binding:"min=0"`

	RoundType              int    `json:"round_type" binding:"max=3,min=0"`
	IpList                 string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList             string `json:"weight_list" binding:"required,valid_weightlist"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" binding:"min=0"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" binding:"min=0"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" binding:"min=0"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" binding:"min=0"`
}

type AddGRPCService struct {
	ServiceName       string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" binding:"required"`
	Port              int    `json:"port" binding:"required,min=8001,max=8999"`
	HeaderTransfer    string `json:"header_transfer" binding:"valid_header_transfer"`
	EnableAuth        int    `json:"enable_auth"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `json:"service_flow_limit"`
	RoundType         int    `json:"round_type"`
	IpList            string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" binding:"required,valid_weightlist"`
}

type UpdateGRPCService struct {
	ID                int    `json:"id" validate:"required"`
	ServiceName       string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" binding:"required"`
	Port              int    `json:"port" binding:"required,min=8001,max=8999"`
	HeaderTransfer    string `json:"header_transfer" binding:"valid_header_transfer"`
	EnableAuth        int    `json:"enable_auth"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `json:"service_flow_limit"`
	RoundType         int    `json:"round_type"`
	IpList            string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" binding:"required,valid_weightlist"`
}

type AddTCPService struct {
	ServiceName       string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" binding:"required"`
	Port              int    `json:"port" binding:"required,min=8001,max=8999"`
	EnableAuth        int    `json:"enable_auth"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `json:"service_flow_limit"`
	RoundType         int    `json:"round_type"`
	IpList            string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" binding:"required,valid_weightlist"`
}

type UpdateTCPService struct {
	ID                int    `json:"id" validate:"required"`
	ServiceName       string `json:"service_name" binding:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" binding:"required"`
	Port              int    `json:"port" binding:"required,min=8001,max=8999"`
	EnableAuth        int    `json:"enable_auth"`
	BlackList         string `json:"black_list"`
	WhiteList         string `json:"white_list"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `json:"service_flow_limit"`
	RoundType         int    `json:"round_type"`
	IpList            string `json:"ip_list" binding:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" binding:"required,valid_weightlist"`
}

func (service *ListService) ServiceList() serializer.Response {
	var serviceInfo model.ServiceInfo
	list, total, err := serviceInfo.PageList(global.DB, service.Info, service.PageNo, service.PageSize)
	if err != nil {
		return serializer.BuildResponseErr(2000, err)
	}

	items := make([]*serializer.ServiceListItem, 0, len(list))

	for _, info := range list {
		serviceDetail, err := info.ServiceDetail(global.DB)
		if err != nil {
			return serializer.BuildResponseErr(2001, err)
		}
		serviceAddr := "unknown"

		clusterIP := global.ProxyConfig.Common.Addr
		clusterPort := global.ProxyConfig.Http.Port
		clusterSslPort := global.ProxyConfig.Https.Port

		switch serviceDetail.Info.LoadType {
		case constants.LoadTypeHTTP:
			switch serviceDetail.HTTPRule.RuleType {
			case constants.HTTPRuleTypePrefixURL: // URL前缀
				if serviceDetail.HTTPRule.NeedHttps == 1 {
					serviceAddr = fmt.Sprintf("%s:%d%s", clusterIP, clusterSslPort, serviceDetail.HTTPRule.Rule)
				} else {
					serviceAddr = fmt.Sprintf("%s:%d%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
				}
			case constants.HTTPRuleTypeDomain: // 域名
				serviceAddr = serviceDetail.HTTPRule.Rule
			}
		case constants.LoadTypeTCP:
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		case constants.LoadTypeGRPC:
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}

		ipList := serviceDetail.LoadBalance.GetIPList()

		serviceCounter, err := handler.FlowCounterHandler.GetCounter(constants.FlowServicePrefix + info.ServiceName)
		if err != nil {
			return serializer.BuildResponseErr(2002, err)
		}
		items = append(items, &serializer.ServiceListItem{
			ID:          info.ID,
			ServiceName: info.ServiceName,
			ServiceDesc: info.ServiceDesc,
			LoadType:    info.LoadType,
			ServiceAddr: serviceAddr,
			Qpd:         serviceCounter.TotalCount,
			Qps:         serviceCounter.QPS,
			TotalNode:   len(ipList),
		})
	}

	return serializer.BuildResponseOkWithData(codes.Success, serializer.ServiceList{
		Total: total,
		List:  items,
	})
}

func (service *DetailService) ServiceDetail() serializer.Response {
	serviceInfo := &model.ServiceInfo{
		ID: uint(service.ID),
	}
	if err := serviceInfo.Find(global.DB); err != nil {
		return serializer.BuildResponseErr(2000, err)
	}
	// 获取服务详细信息
	serviceDetail, err := serviceInfo.ServiceDetail(global.DB)
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOkWithData(codes.Success, serviceDetail)
}

func (service *DeleteService) ServiceDelete() serializer.Response {
	serviceInfo := model.ServiceInfo{
		ID: uint(service.ID),
	}
	if err := serviceInfo.Find(global.DB); err != nil {
		return serializer.BuildResponseErr(2000, err)
	}
	// 标记为已删除
	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(global.DB); err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *AddHTTPService) AddHTTPService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ServiceName: service.ServiceName,
		IsDelete:    0,
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(map[string]interface{}{
			"service_name": serviceInfo.ServiceName,
			"is_delete":    serviceInfo.IsDelete,
		}).First(&serviceInfo).Error; err == nil {
			return errors.New("服务已存在")
		}
		httpRule := model.HttpRule{RuleType: service.RuleType, Rule: service.Rule}
		if err := httpRule.Find(tx); err == nil {
			return errors.New("服务接入前缀或域名已存在")
		}

		// 保存 ServiceInfo 基本信息
		serviceInfo = model.ServiceInfo{
			LoadType:    constants.LoadTypeHTTP, // HTTP
			ServiceName: service.ServiceName,
			ServiceDesc: service.ServiceDesc,
		}

		if err := serviceInfo.Save(tx); err != nil {
			return err
		}

		// 保存 HTTP 规则
		httpRule = model.HttpRule{
			ServiceID:      serviceInfo.ID,
			RuleType:       service.RuleType,
			Rule:           service.Rule,
			NeedHttps:      service.NeedHttps,
			NeedStripUri:   service.NeedStripUri,
			HeaderTransfer: service.HeaderTransfer,
			UrlRewrite:     service.UrlRewrite,
		}
		if err := httpRule.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := model.AccessControl{
			ServiceID:         serviceInfo.ID,
			EnableAuth:        service.EnableAuth,
			BlackList:         service.BlackList,
			WhiteList:         service.WhiteList,
			ClientIPFlowLimit: service.ClientIPFlowLimit,
			ServiceFlowLimit:  service.ServiceFlowLimit,
		}
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := model.LoadBalance{
			ServiceID:              serviceInfo.ID,
			RoundType:              service.RoundType,
			IpList:                 service.IpList,
			WeightList:             service.WeightList,
			UpstreamConnectTimeout: service.UpstreamConnectTimeout,
			UpstreamHeaderTimeout:  service.UpstreamHeaderTimeout,
			UpstreamIdleTimeout:    service.UpstreamIdleTimeout,
			UpstreamMaxIdle:        service.UpstreamMaxIdle,
		}
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *UpdateHTTPService) UpdateHTTPService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ID:          uint(service.ID),
		ServiceName: service.ServiceName,
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := serviceInfo.Find(tx); err != nil {
			return errors.New("服务不存在")
		}

		serviceDetail, err := serviceInfo.ServiceDetail(tx)
		if err != nil {
			return errors.New("服务不存在")
		}

		// 保存 ServiceInfo 基本信息
		info := serviceDetail.Info
		info.ServiceDesc = service.ServiceDesc
		if err := info.Save(tx); err != nil {
			return err
		}

		// 保存 HTTP 规则
		httpRule := serviceDetail.HTTPRule
		httpRule.NeedHttps = service.NeedHttps
		httpRule.NeedStripUri = service.NeedStripUri
		httpRule.UrlRewrite = service.UrlRewrite
		httpRule.HeaderTransfer = service.HeaderTransfer
		if err := httpRule.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := serviceDetail.AccessControl
		accessControl.EnableAuth = service.EnableAuth
		accessControl.BlackList = service.BlackList
		accessControl.WhiteList = service.WhiteList
		accessControl.ClientIPFlowLimit = service.ClientIPFlowLimit
		accessControl.ServiceFlowLimit = service.ServiceFlowLimit
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := serviceDetail.LoadBalance
		loadBalance.RoundType = service.RoundType
		loadBalance.IpList = service.IpList
		loadBalance.WeightList = service.WeightList
		loadBalance.UpstreamConnectTimeout = service.UpstreamConnectTimeout
		loadBalance.UpstreamHeaderTimeout = service.UpstreamHeaderTimeout
		loadBalance.UpstreamIdleTimeout = service.UpstreamIdleTimeout
		loadBalance.UpstreamMaxIdle = service.UpstreamMaxIdle
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *AddGRPCService) AddGRPCService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ServiceName: service.ServiceName,
		IsDelete:    0,
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(map[string]interface{}{
			"service_name": serviceInfo.ServiceName,
			"is_delete":    serviceInfo.IsDelete,
		}).First(&serviceInfo).Error; err == nil {
			return errors.New("服务名被占用，请重新输入")
		}

		// 验证端口是否被TCP和GRPC服务占用
		tcpRule := model.TcpRule{Port: service.Port}
		if err := tcpRule.Find(tx); err == nil {
			return errors.New("服务端口被TCP服务占用，请重新输入")
		}
		grpcRule := model.GrpcRule{Port: service.Port}
		if err := grpcRule.Find(tx); err == nil {
			return errors.New("服务端口被GRPC服务占用，请重新输入")
		}

		// 保存 ServiceInfo 基本信息
		serviceInfo = model.ServiceInfo{
			LoadType:    constants.LoadTypeGRPC, // GRPC
			ServiceName: service.ServiceName,
			ServiceDesc: service.ServiceDesc,
		}

		if err := serviceInfo.Save(tx); err != nil {
			return err
		}

		// 保存 GRPC 规则
		grpcRule = model.GrpcRule{
			ServiceID:      serviceInfo.ID,
			Port:           service.Port,
			HeaderTransfer: service.HeaderTransfer,
		}
		if err := grpcRule.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := model.AccessControl{
			ServiceID:         serviceInfo.ID,
			EnableAuth:        service.EnableAuth,
			BlackList:         service.BlackList,
			WhiteList:         service.WhiteList,
			ClientIPFlowLimit: service.ClientIPFlowLimit,
			ServiceFlowLimit:  service.ServiceFlowLimit,
		}
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := model.LoadBalance{
			ServiceID:  serviceInfo.ID,
			RoundType:  service.RoundType,
			IpList:     service.IpList,
			WeightList: service.WeightList,
		}
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *UpdateGRPCService) UpdateGRPCService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ID:          uint(service.ID),
		ServiceName: service.ServiceName,
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := serviceInfo.Find(tx); err != nil {
			return errors.New("服务不存在")
		}

		serviceDetail, err := serviceInfo.ServiceDetail(tx)
		if err != nil {
			return errors.New("服务不存在")
		}

		// 保存 ServiceInfo 基本信息，服务名不可改
		info := serviceDetail.Info
		info.ServiceDesc = service.ServiceDesc
		if err := info.Save(tx); err != nil {
			return err
		}

		// 保存 GRPC 规则, 端口不可改
		grpcRule := serviceDetail.GRPCRule
		grpcRule.HeaderTransfer = service.HeaderTransfer
		if err := grpcRule.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := serviceDetail.AccessControl
		accessControl.EnableAuth = service.EnableAuth
		accessControl.BlackList = service.BlackList
		accessControl.WhiteList = service.WhiteList
		accessControl.ClientIPFlowLimit = service.ClientIPFlowLimit
		accessControl.ServiceFlowLimit = service.ServiceFlowLimit
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := serviceDetail.LoadBalance
		loadBalance.RoundType = service.RoundType
		loadBalance.IpList = service.IpList
		loadBalance.WeightList = service.WeightList
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *AddTCPService) AddTCPService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ServiceName: service.ServiceName,
		IsDelete:    0,
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(map[string]interface{}{
			"service_name": serviceInfo.ServiceName,
			"is_delete":    serviceInfo.IsDelete,
		}).First(&serviceInfo).Error; err == nil {
			return errors.New("服务名被占用，请重新输入")
		}

		// 验证端口是否被TCP和GRPC服务占用
		tcpRule := model.TcpRule{Port: service.Port}
		if err := tcpRule.Find(tx); err == nil {
			return errors.New("服务端口被TCP服务占用，请重新输入")
		}
		grpcRule := model.GrpcRule{Port: service.Port}
		if err := grpcRule.Find(tx); err == nil {
			return errors.New("服务端口被GRPC服务占用，请重新输入")
		}

		// 保存 ServiceInfo 基本信息
		serviceInfo = model.ServiceInfo{
			LoadType:    constants.LoadTypeTCP, // TCP
			ServiceName: service.ServiceName,
			ServiceDesc: service.ServiceDesc,
		}

		if err := serviceInfo.Save(tx); err != nil {
			return err
		}

		// 保存 TCP 规则
		tcpRule = model.TcpRule{
			ServiceID: serviceInfo.ID,
			Port:      service.Port,
		}
		if err := tcpRule.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := model.AccessControl{
			ServiceID:         serviceInfo.ID,
			EnableAuth:        service.EnableAuth,
			BlackList:         service.BlackList,
			WhiteList:         service.WhiteList,
			ClientIPFlowLimit: service.ClientIPFlowLimit,
			ServiceFlowLimit:  service.ServiceFlowLimit,
		}
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := model.LoadBalance{
			ServiceID:  serviceInfo.ID,
			RoundType:  service.RoundType,
			IpList:     service.IpList,
			WeightList: service.WeightList,
		}
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}

func (service *UpdateTCPService) UpdateTCPService() serializer.Response {
	if len(strings.Split(service.IpList, ",")) !=
		len(strings.Split(service.WeightList, ",")) {
		return serializer.BuildResponseErr(2000, errors.New("IP列表与权重列表数量不一致"))
	}

	serviceInfo := model.ServiceInfo{
		ID:          uint(service.ID),
		ServiceName: service.ServiceName,
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := serviceInfo.Find(tx); err != nil {
			return errors.New("服务不存在")
		}

		serviceDetail, err := serviceInfo.ServiceDetail(tx)
		if err != nil {
			return errors.New("服务不存在")
		}

		// 保存 ServiceInfo 基本信息，服务名不可改
		info := serviceDetail.Info
		info.ServiceDesc = service.ServiceDesc
		if err := info.Save(tx); err != nil {
			return err
		}

		// 保存 访问控制规则
		accessControl := serviceDetail.AccessControl
		accessControl.EnableAuth = service.EnableAuth
		accessControl.BlackList = service.BlackList
		accessControl.WhiteList = service.WhiteList
		accessControl.ClientIPFlowLimit = service.ClientIPFlowLimit
		accessControl.ServiceFlowLimit = service.ServiceFlowLimit
		if err := accessControl.Save(tx); err != nil {
			return err
		}

		// 保存 负载均衡规则
		loadBalance := serviceDetail.LoadBalance
		loadBalance.RoundType = service.RoundType
		loadBalance.IpList = service.IpList
		loadBalance.WeightList = service.WeightList
		if err := loadBalance.Save(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return serializer.BuildResponseErr(2001, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}
