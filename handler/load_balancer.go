package handler

import (
	"fmt"

	"github.com/josexy/gw/global"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/pkg/etcd"
	"github.com/josexy/gw/reverse_proxy/load_balance"
)

var LoadBalancerHandler *LoadBalancer

type LoadBalancer struct {
	handler MapBaseHandler[load_balance.LoadBalance]
}

func init() {
	LoadBalancerHandler = NewLoadBalancer()
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		handler: NewMapBaseHandler[load_balance.LoadBalance](),
	}
}

func (lbr *LoadBalancer) GetLoadBalancer(serviceDetail *model.ServiceDetail) (load_balance.LoadBalance, error) {
	lbr.handler.RLock()
	if lb, ok := lbr.handler.Cache[serviceDetail.Info.ServiceName]; ok {
		lbr.handler.RUnlock()
		return lb, nil
	}
	lbr.handler.RUnlock()

	schema := "http://"
	if serviceDetail.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}
	if serviceDetail.Info.LoadType == constants.LoadTypeTCP ||
		serviceDetail.Info.LoadType == constants.LoadTypeGRPC {
		schema = ""
	}

	// 上游节点地址
	ipList := serviceDetail.LoadBalance.GetIPList()
	weightList := serviceDetail.LoadBalance.GetWeightList()

	// http://127.0.0.1:8080 -> 50
	// 127.0.0.2:2221 -> 20
	ipWeightList := make(map[string]string, len(ipList))
	for index, item := range ipList {
		ipWeightList[item] = weightList[index]
	}

	discovery, err := etcd.NewServiceDiscovery(global.AppConfig.Etcd.Endpoints,
		fmt.Sprintf("%s%s", schema, "%s"), ipWeightList)

	if err != nil {
		return nil, err
	}

	// 服务监控
	serviceList, err := discovery.WatchService(constants.EtcdPrefix)
	if err != nil {
		return nil, err
	}

	lb := load_balance.NewLoadBalanceFromFactorWithDiscovery(
		load_balance.LbType(serviceDetail.LoadBalance.RoundType), discovery)

	// 初始化负载均衡器服务地址
	if obs, ok := lb.(etcd.Observer); ok {
		obs.Update(serviceList)
	}

	lbr.handler.List = append(lbr.handler.List, lb)
	lbr.handler.Lock()
	lbr.handler.Cache[serviceDetail.Info.ServiceName] = lb
	lbr.handler.Unlock()
	return lb, nil
}
