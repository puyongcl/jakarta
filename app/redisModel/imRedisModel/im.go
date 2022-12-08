package imRedisModel

import (
	"context"
	"jakarta/common/key/rediskey"
	"strconv"
)

// 多少分钟内不回复消息自动转状态为休息中
// 放入 不存在 则记录 存在则不更新 并返回记录的时间
func (m *IMRedis) MarkUserReplyTime(ctx context.Context, listenerUid int64, stamp string) (string, error) {
	b, err := m.redisClient.HsetnxCtx(ctx, rediskey.RedisKeyMarkReplyMsgTime, strconv.FormatInt(listenerUid, 10), stamp)
	if err != nil {
		return "", err
	}
	if b {
		return stamp, nil
	}
	return m.redisClient.HgetCtx(ctx, rediskey.RedisKeyMarkReplyMsgTime, strconv.FormatInt(listenerUid, 10))
}

// 移除
func (m *IMRedis) CancelMarkUserReplyTime(ctx context.Context, listenerUid string) error {
	_, err := m.redisClient.HdelCtx(ctx, rediskey.RedisKeyMarkReplyMsgTime, listenerUid)
	return err
}

//
func (m *IMRedis) ScanMarkUserReplyTime(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error) {
	return m.redisClient.HscanCtx(ctx, rediskey.RedisKeyMarkReplyMsgTime, cursor, "*", count)
}
