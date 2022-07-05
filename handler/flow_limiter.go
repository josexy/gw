package handler

import (
	"golang.org/x/time/rate"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	handler MapBaseHandler[*rate.Limiter]
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		handler: NewMapBaseHandler[*rate.Limiter](),
	}
}

func (fc *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	fc.handler.RLock()
	if limiter, ok := fc.handler.Cache[serviceName]; ok {
		fc.handler.RUnlock()
		return limiter, nil
	}
	fc.handler.RUnlock()

	// 令牌桶限流器
	limiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))

	fc.handler.Lock()
	fc.handler.List = append(fc.handler.List, limiter)
	fc.handler.Cache[serviceName] = limiter
	fc.handler.Unlock()
	return limiter, nil
}
