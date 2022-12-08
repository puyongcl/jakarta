package userRedisModel

import (
	"context"
	"jakarta/common/key/rediskey"
)

func (r *UserRedis) IncrUidIdx(ctx context.Context) (res int64, err error) {
	return r.redisClient.IncrCtx(ctx, rediskey.RedisKeyUidIndex)
}
