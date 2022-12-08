package statisticRedisModel

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"jakarta/common/dashboard"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"jakarta/common/tool"
	"strconv"
	"time"
)

// 初始化XXX看板数据
func (m *StatisticRedis) InitListenerDashboard(ctx context.Context, listenerUid int64) error {
	fields := tool.GetStructFieldName(dashboard.ListenerDashboardRedisHashData{})
	fm := make(map[string]string, 0)
	for idx := 0; idx < len(fields); idx++ {
		fm[fields[idx]] = "0"
	}
	fm["ListenerUid"] = strconv.FormatInt(listenerUid, 10)
	fm["LastDayStatUpdateTime"] = time.Now().Format(db.DateTimeFormat)
	fm["TodayStatUpdateTime"] = time.Now().Format(db.DateTimeFormat)
	return m.redisClient.HmsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid), fm)
}

// 更新XXX看板数据
func (m *StatisticRedis) HSetListenerDashboard(ctx context.Context, listenerUid int64, field string, value int64) error {
	return m.redisClient.HsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid), field, strconv.FormatInt(value, 10))
}
func (m *StatisticRedis) HMSetListenerDashboard(ctx context.Context, listenerUid int64, fm map[string]string) error {
	return m.redisClient.HmsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid), fm)
}

// 获取XXX看板数据
func (m *StatisticRedis) GetListenerDashboard(ctx context.Context, listenerUid int64) (*dashboard.ListenerDashboardRedisHashData, error) {
	res, err := m.redisClient.HgetallCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid))
	if err != nil {
		return nil, err
	}
	dsb := new(dashboard.ListenerDashboardRedisHashData)
	err = mapstructure.Decode(res, dsb)
	if err != nil {
		return nil, err
	}
	return dsb, nil
}
