package listenerRedisModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"strconv"
	"time"
)

// 推荐XXX 记录推荐次数 每次推荐最小值

// 批量加入列表
func (m *ListenerRedis) InitRecommendListenerPool(ctx context.Context, keyFmt string, l []int64) error {
	return m.redisClient.PipelinedCtx(ctx, func(p redis.Pipeliner) error {
		var mb []*redis.Z

		for k, _ := range l {
			mb = append(mb, &redis.Z{Score: 0, Member: l[k]})
		}

		err := p.ZAddNX(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), mb...).Err()
		return err
	})
}

// 新加一个 获取当前最小score 作为初始值 TODO 改为批量接口
func (m *ListenerRedis) ADDNewUserRecommendListenerOne(ctx context.Context, keyFmt string, listenerUid int64) error {
	r, err := m.redisClient.ZrangebyscoreWithScoresAndLimitCtx(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), 0, listenerkey.MaxNewUserRecommendScore, 0, 1)
	if err != nil {
		return err
	}
	var score int64
	if len(r) >= 1 {
		score = r[0].Score
	}

	_, err = m.redisClient.ZaddCtx(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), score, strconv.FormatInt(listenerUid, 10))
	return err
}

// 推荐几个 分值最小的并加1
func (m *ListenerRedis) GetRecommendListenerAndIncrScore(ctx context.Context, keyFmt string, cnt int) ([]int64, error) {
	r, err := m.redisClient.ZrangebyscoreWithScoresAndLimitCtx(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), 0, listenerkey.MaxNewUserRecommendScore, 0, cnt)
	if err != nil {
		return []int64{}, err
	}

	var rv []int64
	for idx := 0; idx < len(r); idx++ {
		var v int64
		v, err = strconv.ParseInt(r[idx].Key, 10, 64)
		if err != nil {
			return rv, err
		}
		rv = append(rv, v)
	}

	if len(rv) <= 0 {
		return rv, nil
	}

	// 加1
	err = m.redisClient.PipelinedCtx(ctx, func(p redis.Pipeliner) error {
		for k, _ := range rv {
			err2 := p.ZIncrBy(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), 1, strconv.FormatInt(rv[k], 10)).Err()
			if err2 != nil {
				return err2
			}
		}
		return nil
	})

	return rv, err
}

// 移除推荐列表 TODO 改为批量接口
func (m *ListenerRedis) RemoveNewUserRecommendListenerOne(ctx context.Context, keyFmt string, listenerUid int64) error {
	_, err := m.redisClient.ZremCtx(ctx, fmt.Sprintf(keyFmt, time.Now().Format(db.DateFormat2)), strconv.FormatInt(listenerUid, 10))
	return err
}
