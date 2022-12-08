package listenerRedisModel

import "github.com/zeromicro/go-zero/core/stores/redis"

type ListenerRedis struct {
	redisClient *redis.Redis
}

// NewListenerRedis .
func NewListenerRedis(rc *redis.Redis) *ListenerRedis {
	return &ListenerRedis{
		redisClient: rc,
	}
}
