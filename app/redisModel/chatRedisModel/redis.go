package chatRedisModel

import "github.com/zeromicro/go-zero/core/stores/redis"

type ChatRedis struct {
	redisClient *redis.Redis
}

// NewChatRedis
func NewChatRedis(rc *redis.Redis) *ChatRedis {
	return &ChatRedis{
		redisClient: rc,
	}
}
