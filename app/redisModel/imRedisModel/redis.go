package imRedisModel

import "github.com/zeromicro/go-zero/core/stores/redis"

type IMRedis struct {
	redisClient *redis.Redis
}

// NewIMRedis .
func NewIMRedis(rc *redis.Redis) *IMRedis {
	return &IMRedis{
		redisClient: rc,
	}
}
