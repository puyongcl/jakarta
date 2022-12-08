package imRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"time"
)

func (m *IMRedis) IsHaveNotifyLimit(ctx context.Context, id string, msgType int64) (bool, error) {
	v, err := m.redisClient.GetCtx(ctx, fmt.Sprintf(rediskey.RedisKeyNotifyMsgLimit, id, msgType))
	if err != nil {
		return false, err
	}
	if v != "" {
		return true, nil
	}
	return false, nil
}

func (m *IMRedis) SetNotifyLimit(ctx context.Context, id string, msgType int64, sec int) (bool, error) {
	if sec <= 0 { // 不限制
		return true, nil
	}
	return m.redisClient.SetnxExCtx(ctx, fmt.Sprintf(rediskey.RedisKeyNotifyMsgLimit, id, msgType), time.Now().Format(db.DateTimeFormat), sec)
}

func (m *IMRedis) RenewNotifyLimit(ctx context.Context, id string, msgType int64, sec int) error {
	if sec <= 0 { // 不限制
		return nil
	}
	return m.redisClient.SetexCtx(ctx, fmt.Sprintf(rediskey.RedisKeyNotifyMsgLimit, id, msgType), time.Now().Format(db.DateTimeFormat), sec)
}
