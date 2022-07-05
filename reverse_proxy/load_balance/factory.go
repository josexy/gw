package load_balance

import (
	"github.com/josexy/gw/pkg/etcd"
)

func NewLoadBalanceFromFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbConsistentHash:
		return NewConsistentHashBalance(10, nil)
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	default:
		return &RandomBalance{}
	}
}

func NewLoadBalanceFromFactorWithDiscovery(lbType LbType, discovery *etcd.ServiceDiscovery) LoadBalance {
	lb := NewLoadBalanceFromFactory(lbType)
	discovery.AddObserver(lb.(etcd.Observer))
	return lb
}
