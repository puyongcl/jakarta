package redisModel

import "github.com/zeromicro/go-zero/core/stores/redis"

type JobRedis struct {
	redisClient *redis.Redis
}

// NewJobRedis
func NewJobRedis(rc *redis.Redis) *JobRedis {
	return &JobRedis{
		redisClient: rc,
	}
}
