package imRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/rediskey"
	"strconv"
)

// 添加到订阅消息成员中
func (m *IMRedis) AddSubscribeNotifyMember(ctx context.Context, targetUid, msgType, sendCnt, uid int64) error {
	_, err := m.redisClient.SaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeNotifyMember, targetUid, msgType, sendCnt), uid)
	return err
}

// 获取消息订阅成员
func (m *IMRedis) GetSubscribeNotifyMember(ctx context.Context, targetUid, msgType, sendCnt int64, cnt int) ([]string, error) {
	return m.redisClient.SrandmemberCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeNotifyMember, targetUid, msgType, sendCnt), cnt)
}

// 移除订阅消息成员
func (m *IMRedis) RemSubscribeNotifyMember(ctx context.Context, targetUid, msgType, sendCnt int64, val ...string) (int, error) {
	var r []interface{}
	for k, _ := range val {
		r = append(r, val[k])
	}
	return m.redisClient.SremCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeNotifyMember, targetUid, msgType, sendCnt), r...)
}

// 一次性订阅消息
func (m *IMRedis) AddOneTimeSubscribeMsgMember(ctx context.Context, msgType, uid int64) error {
	return m.redisClient.HsetCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeOneTimeNotifyMsgMember, msgType), strconv.FormatInt(uid, 10), "1")
}

// 获取是否有订阅
func (m *IMRedis) IsUserHaveOneTimeSubscribeMsg(ctx context.Context, msgType, uid int64) (bool, error) {
	return m.redisClient.HexistsCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeOneTimeNotifyMsgMember, msgType), strconv.FormatInt(uid, 10))
}

// 发送消息 减1 并删除field
func (m *IMRedis) CostOneTimeSubscribeMsg(ctx context.Context, msgType, uid int64) (bool, error) {
	return m.redisClient.HdelCtx(ctx, fmt.Sprintf(rediskey.RedisKeySubscribeOneTimeNotifyMsgMember, msgType), strconv.FormatInt(uid, 10))
}
