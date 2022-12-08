package chatRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"strconv"
	"time"
)

// 今日进入XXX聊天页面用户数 zset XXX 用户数 (获取数据定义在 listener服务中)
// 设置今日进入XXX聊天页面用户数
func (r *ChatRedis) SetTodayListenerEnterChatUserCnt(ctx context.Context, listenerUid int64, cnt int64) error {
	_, err := r.redisClient.ZaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayEnterChatUserCnt, time.Now().Format(db.DateFormat2)), cnt, strconv.FormatInt(listenerUid, 10))
	return err
}

// 更新XXX看板数据
func (r *ChatRedis) SetListenerDashboard(ctx context.Context, listenerUid int64, field string, value int64) error {
	return r.redisClient.HsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid), field, strconv.FormatInt(value, 10))
}
func (r *ChatRedis) HMSetListenerDashboard(ctx context.Context, listenerUid int64, fm map[string]string) error {
	return r.redisClient.HmsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerDashboard, listenerUid), fm)
}
