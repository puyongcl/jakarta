package redisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"time"
)

// 清除多少天以上的历史每日统计数据
func (r *JobRedis) ClearTodayStatHistoryData(ctx context.Context, keyFmt []string) error {
	now := time.Now()
	tt := now.AddDate(0, 0, -listenerkey.AutoClearListenerTodayStatHistoryData).Format(db.DateFormat2)
	tt2 := now.AddDate(0, 0, -(listenerkey.AutoClearListenerTodayStatHistoryData + 1)).Format(db.DateFormat2)
	tt3 := now.AddDate(0, 0, -(listenerkey.AutoClearListenerTodayStatHistoryData + 2)).Format(db.DateFormat2)
	key := make([]string, len(keyFmt))
	for idx := 0; idx < len(keyFmt); idx++ {
		key = append(key, fmt.Sprintf(keyFmt[idx], tt))
		key = append(key, fmt.Sprintf(keyFmt[idx], tt2))
		key = append(key, fmt.Sprintf(keyFmt[idx], tt3))
	}
	_, err := r.redisClient.DelCtx(ctx, key...)
	return err
}

// 清除多少天以上的每日推荐XXX数据
func (r *JobRedis) ClearRecommendListenerPoolData(ctx context.Context, keyFmt []string) error {
	now := time.Now()
	tt := now.AddDate(0, 0, -listenerkey.AutoClearNewUserRecommendListenerData).Format(db.DateFormat2)
	tt2 := now.AddDate(0, 0, -(listenerkey.AutoClearNewUserRecommendListenerData + 1)).Format(db.DateFormat2)
	tt3 := now.AddDate(0, 0, -(listenerkey.AutoClearNewUserRecommendListenerData + 2)).Format(db.DateFormat2)
	key := make([]string, len(keyFmt))
	for idx := 0; idx < len(keyFmt); idx++ {
		key = append(key, fmt.Sprintf(keyFmt[idx], tt))
		key = append(key, fmt.Sprintf(keyFmt[idx], tt2))
		key = append(key, fmt.Sprintf(keyFmt[idx], tt3))
	}
	_, err := r.redisClient.DelCtx(ctx, key...)
	return err
}
