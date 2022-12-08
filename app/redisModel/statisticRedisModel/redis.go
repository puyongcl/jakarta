package statisticRedisModel

import "github.com/zeromicro/go-zero/core/stores/redis"

type StatisticRedis struct {
	redisClient *redis.Redis
}

// NewStatisticRedis .
func NewStatisticRedis(rc *redis.Redis) *StatisticRedis {
	return &StatisticRedis{
		redisClient: rc,
	}
}
