package userRedisModel

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type UserRedis struct {
	redisClient *redis.Redis
}

// NewUserRedis .
func NewUserRedis(rc *redis.Redis) *UserRedis {
	return &UserRedis{
		redisClient: rc,
	}
}
