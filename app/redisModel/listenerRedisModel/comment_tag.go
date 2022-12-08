package listenerRedisModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/common/key/listenerkey"
	"strconv"
)

// XXX评价标签 有序集合 存放 评价列表上方需要展示 评价标签统计情况

func (m *ListenerRedis) AddCommentTag(ctx context.Context, key string, listenerUid int64, tag []int64) error {
	if len(tag) <= 0 {
		return nil
	}
	err := m.redisClient.PipelinedCtx(ctx, func(pip redis.Pipeliner) error {
		for k, _ := range tag {
			err := pip.ZIncrBy(ctx, fmt.Sprintf(key, listenerUid), 1, strconv.FormatInt(tag[k], 10)).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *ListenerRedis) GetTopCommentTagStat(ctx context.Context, key string, listenerUid int64) ([]redis.Pair, error) {
	return m.redisClient.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, fmt.Sprintf(key, listenerUid), 0, 999999, 0, listenerkey.ShowTopCommentTagCnt)
}

func (m *ListenerRedis) DelCommentTagStat(ctx context.Context, key string, listenerUid int64) error {
	_, err := m.redisClient.DelCtx(ctx, fmt.Sprintf(key, listenerUid))
	return err
}
