package chatRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/tool"
	"strconv"
)

// 进入聊天界面 首次开始聊天
func (r *ChatRedis) GetSetListenerFreeMsgCnt(ctx context.Context, uid, listenerUid int64) (int64, error, bool) {
	key := fmt.Sprintf(rediskey.RedisKeyUserPerListenerFreeTextMsgCnt, uid, listenerUid)
	rs, err := r.redisClient.GetCtx(ctx, key)
	if err != nil {
		return 0, err, false
	}
	if rs == "" {
		err = r.redisClient.SetexCtx(ctx, key, strconv.FormatInt(chatkey.FreeTextChatCntPerDay, 10), int(tool.GetTodayRemainSecond()))
		if err != nil {
			return 0, err, false
		}
		return chatkey.FreeTextChatCntPerDay, nil, true
	}
	cnt, err := strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return 0, err, false
	}
	return cnt, nil, false
}

func (r *ChatRedis) ResetListenerFreeMsgCnt(ctx context.Context, uid, listenerUid int64) error {
	key := fmt.Sprintf(rediskey.RedisKeyUserPerListenerFreeTextMsgCnt, uid, listenerUid)
	return r.redisClient.SetexCtx(ctx, key, strconv.FormatInt(chatkey.FreeTextChatCntPerDay, 10), int(tool.GetTodayRemainSecond()))
}

// 聊天过程免费聊天次数使用
func (r *ChatRedis) DecrListenerFreeMsgCnt(ctx context.Context, uid, listenerUid int64) (int64, error) {
	key := fmt.Sprintf(rediskey.RedisKeyUserPerListenerFreeTextMsgCnt, uid, listenerUid)
	cnt, err := r.redisClient.DecrCtx(ctx, key)
	if cnt < 0 {
		cnt = 0
	}
	return cnt, err
}

// 下语音订单赠送免费聊天次数
func (r *ChatRedis) AddListenerFreeChatCnt(ctx context.Context, uid, listenerUid int64) (int64, error) {
	key := fmt.Sprintf(rediskey.RedisKeyUserPerListenerFreeTextMsgCnt, uid, listenerUid)
	rs, err := r.redisClient.GetCtx(ctx, key)
	if err != nil {
		return 0, err
	}
	var cnt int64
	if rs == "" {
		cnt = int64(chatkey.FreeTextChatCntPerDay + chatkey.GiftTextChatCnt)
		err = r.redisClient.SetexCtx(ctx, key, strconv.FormatInt(cnt, 10), int(tool.GetTodayRemainSecond()))
		if err != nil {
			return 0, err
		}
		return cnt, nil
	}
	cnt, err = strconv.ParseInt(rs, 10, 64)
	if err != nil {
		return 0, err
	}
	if cnt <= 0 {
		cnt = chatkey.GiftTextChatCnt + tool.Abs(cnt)
	}
	cnt, err = r.redisClient.IncrbyCtx(ctx, key, cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
