package handler

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/josexy/gw/global"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/pkg/constants"
)

// RedisFlowCountService 流量统计
type RedisFlowCountService struct {
	ServiceName string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(serviceName string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		ServiceName: serviceName,
		Interval:    interval,
		QPS:         0,
		Unix:        0,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Debug("%v", err)
			}
		}()
		// 统计一秒内用户访问的总请求量和QPS
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C

			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			currentTime := time.Now()

			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)

			_, err := global.Redis.Pipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
				pipeliner.IncrBy(context.Background(), dayKey, tickerCount)
				pipeliner.Expire(context.Background(), dayKey, 2*24*time.Hour)
				pipeliner.IncrBy(context.Background(), hourKey, tickerCount)
				pipeliner.Expire(context.Background(), hourKey, 2*24*time.Hour)
				return nil
			})
			if err != nil {
				logx.Debug("redis flow count service err: %v", err)
				continue
			}

			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				continue
			}
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}
			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				// 当日请求量
				reqCounter.TotalCount = totalCount
				// 当前QPS
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (fcs *RedisFlowCountService) GetDayKey(t time.Time) string {
	// flow_day_count_20220603_abc
	return fmt.Sprintf("%s_%s_%s", constants.RedisFlowDayKey, t.Format("20060102"), fcs.ServiceName)
}

func (fcs *RedisFlowCountService) GetHourKey(t time.Time) string {
	// flow_hour_count_2022060315_abc
	return fmt.Sprintf("%s_%s_%s", constants.RedisFlowHourKey, t.Format("2006010215"), fcs.ServiceName)
}

func (fcs *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return global.Redis.Get(context.Background(), fcs.GetHourKey(t)).Int64()
}

func (fcs *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return global.Redis.Get(context.Background(), fcs.GetDayKey(t)).Int64()
}

func (fcs *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Debug("RedisFlowCountService increase err: %v", err)
			}
		}()
		atomic.AddInt64(&fcs.TickerCount, 1)
	}()
}
