package statisticRedisModel

import (
	"context"
	"encoding/json"
	"jakarta/common/average"
	"jakarta/common/key/rediskey"
)

// 一段时间内的昨日平均值
func (m *StatisticRedis) SetListenerAverageStat(ctx context.Context, a *average.ListenerStatAverage) (err error) {
	var buf []byte
	buf, err = json.Marshal(a)
	if err != nil {
		return err
	}
	return m.redisClient.SetCtx(ctx, rediskey.RedisKeyListenerLastDayAverage, string(buf))
}

func (m *StatisticRedis) GetListenerAverageStat(ctx context.Context) (a *average.ListenerStatAverage, err error) {
	v, err := m.redisClient.GetCtx(ctx, rediskey.RedisKeyListenerLastDayAverage)
	if err != nil {
		return nil, err
	}
	a = &average.ListenerStatAverage{}
	if v == "" {
		return
	}
	err = json.Unmarshal([]byte(v), a)
	if err != nil {
		return nil, err
	}
	return
}
