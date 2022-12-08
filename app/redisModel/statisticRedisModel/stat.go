package statisticRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"strconv"
	"time"
)

// 今日接单量 zset XXX 接单量
// 增加今日接单量
func (m *StatisticRedis) AddTodayListenerOrderCnt(ctx context.Context, listenerUid int64, addCnt int64) error {
	_, err := m.redisClient.ZincrbyCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayOrderCnt, time.Now().Format(db.DateFormat2)), addCnt, strconv.FormatInt(listenerUid, 10))
	return err
}

// 获取排名和接单量
func (m *StatisticRedis) GetListenerTodayOrderCntAndRank(ctx context.Context, listenerUid int64) (rank, score int64, err error) {
	key := fmt.Sprintf(rediskey.RedisKeyListenerTodayOrderCnt, time.Now().Format(db.DateFormat2))
	v := strconv.FormatInt(listenerUid, 10)
	rank, err = m.redisClient.ZrevrankCtx(ctx, key, v)
	if err != nil {
		return
	}
	score, err = m.redisClient.ZscoreCtx(ctx, key, v)
	if err != nil {
		return
	}
	return
}

// 今日接单金额 zset XXX 接单金额
// 增加今日接单金额
func (m *StatisticRedis) AddTodayListenerOrderAmount(ctx context.Context, listenerUid int64, addAmount int64) error {
	_, err := m.redisClient.ZincrbyCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayOrderAmount, time.Now().Format(db.DateFormat2)), addAmount, strconv.FormatInt(listenerUid, 10))
	return err
}

// 获取排名和接单金额
func (m *StatisticRedis) GetListenerTodayOrderAmountAndRank(ctx context.Context, listenerUid int64) (rank, score int64, err error) {
	key := fmt.Sprintf(rediskey.RedisKeyListenerTodayOrderAmount, time.Now().Format(db.DateFormat2))
	v := strconv.FormatInt(listenerUid, 10)
	rank, err = m.redisClient.ZrevrankCtx(ctx, key, v)
	if err != nil {
		return
	}
	score, err = m.redisClient.ZscoreCtx(ctx, key, v)
	if err != nil {
		return
	}
	return
}

// 今日曝光推荐用户数 zset XXX 用户数
// 设置今日推荐用户数
func (m *StatisticRedis) SetTodayListenerRecommendUserCnt(ctx context.Context, listenerUid int64, cnt int64) error {
	_, err := m.redisClient.ZaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayRecommendUserCnt, time.Now().Format(db.DateFormat2)), cnt, strconv.FormatInt(listenerUid, 10))
	return err
}

// 获取排名和推荐用户数
func (m *StatisticRedis) GetListenerTodayRecommendUserCntAndRank(ctx context.Context, listenerUid int64) (rank, score int64, err error) {
	key := fmt.Sprintf(rediskey.RedisKeyListenerTodayRecommendUserCnt, time.Now().Format(db.DateFormat2))
	v := strconv.FormatInt(listenerUid, 10)
	rank, err = m.redisClient.ZrevrankCtx(ctx, key, v)
	if err != nil {
		return
	}
	score, err = m.redisClient.ZscoreCtx(ctx, key, v)
	if err != nil {
		return
	}
	return
}

// 今日访问个人资料用户数 zset XXX 用户数
// 设置今日访问个人资料用户数
func (m *StatisticRedis) SetTodayListenerViewUserCnt(ctx context.Context, listenerUid int64, cnt int64) error {
	_, err := m.redisClient.ZaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayViewUserCnt, time.Now().Format(db.DateFormat2)), cnt, strconv.FormatInt(listenerUid, 10))
	return err
}

// 获取排名和访问个人资料用户数
func (m *StatisticRedis) GetListenerTodayViewUserCntAndRank(ctx context.Context, listenerUid int64) (rank, score int64, err error) {
	key := fmt.Sprintf(rediskey.RedisKeyListenerTodayViewUserCnt, time.Now().Format(db.DateFormat2))
	v := strconv.FormatInt(listenerUid, 10)
	rank, err = m.redisClient.ZrevrankCtx(ctx, key, v)
	if err != nil {
		return
	}
	score, err = m.redisClient.ZscoreCtx(ctx, key, v)
	if err != nil {
		return
	}
	return
}

// 今日进入XXX聊天页面用户数 zset XXX 用户数
// 设置今日进入XXX聊天页面用户数 (更新数据在 chat服务中)
func (m *StatisticRedis) SetTodayListenerEnterChatUserCnt(ctx context.Context, listenerUid int64, cnt int64) error {
	_, err := m.redisClient.ZaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeyListenerTodayEnterChatUserCnt, time.Now().Format(db.DateFormat2)), cnt, strconv.FormatInt(listenerUid, 10))
	return err
}

// 获取排名和进入XXX聊天页面用户数
func (m *StatisticRedis) GetListenerTodayEnterChatCntAndRank(ctx context.Context, listenerUid int64) (rank, score int64, err error) {
	key := fmt.Sprintf(rediskey.RedisKeyListenerTodayEnterChatUserCnt, time.Now().Format(db.DateFormat2))
	v := strconv.FormatInt(listenerUid, 10)
	rank, err = m.redisClient.ZrevrankCtx(ctx, key, v)
	if err != nil {
		return
	}
	score, err = m.redisClient.ZscoreCtx(ctx, key, v)
	if err != nil {
		return
	}
	return
}
