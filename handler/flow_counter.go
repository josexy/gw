package handler

import (
	"time"
)

var FlowCounterHandler *FlowCounter

type FlowCounter struct {
	handler MapBaseHandler[*RedisFlowCountService]
}

func init() {
	FlowCounterHandler = NewFlowCounter()
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		handler: NewMapBaseHandler[*RedisFlowCountService](),
	}
}

func (fc *FlowCounter) GetCounter(serviceName string) (*RedisFlowCountService, error) {
	fc.handler.RLock()
	if counter, ok := fc.handler.Cache[serviceName]; ok {
		fc.handler.RUnlock()
		return counter, nil
	}
	fc.handler.RUnlock()

	counter := NewRedisFlowCountService(serviceName, time.Second)
	fc.handler.Lock()
	fc.handler.List = append(fc.handler.List, counter)
	fc.handler.Cache[serviceName] = counter
	fc.handler.Unlock()
	return counter, nil
}
