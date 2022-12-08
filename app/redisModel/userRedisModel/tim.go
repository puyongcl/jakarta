package userRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/rediskey"
)

func (r *UserRedis) SetTimUserSignature(ctx context.Context, uid int64, sig string, expire int) error {
	return r.redisClient.SetexCtx(ctx, fmt.Sprintf(rediskey.CacheTIMUserSignKey, uid), sig, expire)
}

func (r *UserRedis) GetTimUserSignature(ctx context.Context, uid int64) (sig string, err error) {
	return r.redisClient.GetCtx(ctx, fmt.Sprintf(rediskey.CacheTIMUserSignKey, uid))
}
