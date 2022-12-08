package userRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/rediskey"
	"strconv"
)

func (r *UserRedis) AddBlacklist(ctx context.Context, uid, targetUid int64) error {
	_, err := r.redisClient.SaddCtx(ctx, fmt.Sprintf(rediskey.RedisKeyUserBlacklist, uid), targetUid)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRedis) DelBlacklist(ctx context.Context, uid, targetUid int64) error {
	_, err := r.redisClient.SremCtx(ctx, fmt.Sprintf(rediskey.RedisKeyUserBlacklist, uid), targetUid)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRedis) CountBlacklist(ctx context.Context, uid int64) (int64, error) {
	cnt, err := r.redisClient.ScardCtx(ctx, fmt.Sprintf(rediskey.RedisKeyUserBlacklist, uid))
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *UserRedis) GetBlacklist(ctx context.Context, uid int64) ([]int64, error) {
	res, err := r.redisClient.SmembersCtx(ctx, fmt.Sprintf(rediskey.RedisKeyUserBlacklist, uid))
	if err != nil {
		return nil, err
	}
	if len(res) <= 0 {
		return []int64{}, err
	}
	var res2 []int64
	var val int64
	for idx := 0; idx < len(res); idx++ {
		val, err = strconv.ParseInt(res[idx], 10, 64)
		if err != nil {
			return nil, err
		}
		res2 = append(res2, val)
	}
	return res2, nil
}
