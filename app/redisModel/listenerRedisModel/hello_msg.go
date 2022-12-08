package listenerRedisModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/common/key/rediskey"
	"strconv"
	"time"
)

// 用户登陆发送问候语
func (m *ListenerRedis) SetUserHelloMsg(ctx context.Context, uid int64, msg []interface{}) error {
	return m.redisClient.PipelinedCtx(ctx, func(pip redis.Pipeliner) error {
		err := pip.SAdd(ctx, fmt.Sprintf(rediskey.RedisKeyUserLoginHelloMsg, uid), msg...).Err()
		if err != nil {
			return err
		}
		err = pip.Expire(ctx, fmt.Sprintf(rediskey.RedisKeyUserLoginHelloMsg, uid), rediskey.RedisKeyUserLoginHelloMsgExpire*time.Second).Err()
		return err
	})
}

// 取出几条问候语
func (m *ListenerRedis) PopUserHelloMsg(ctx context.Context, uid int64, cnt int64) ([]int64, error) {
	var ri []int64
	var rv int64
	var rs string
	var idx int64
	for idx = 0; idx < cnt; idx++ {
		var err error
		rs, err = m.redisClient.SpopCtx(ctx, fmt.Sprintf(rediskey.RedisKeyUserLoginHelloMsg, uid))
		if err != nil {
			return nil, err
		}
		rv, err = strconv.ParseInt(rs, 10, 64)
		if err != nil {
			return nil, err
		}
		ri = append(ri, rv)
	}
	return ri, nil
}
