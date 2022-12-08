package chatRedisModel

import (
	"context"
	"fmt"
	"jakarta/common/key/rediskey"
	"strconv"
)

// 记录是谁锁定了通话对象
func (r *ChatRedis) RecordLockTargetUser(ctx context.Context, fromUid, toUid, sec int64) (bool, error) {
	return r.redisClient.SetnxExCtx(ctx, fmt.Sprintf(rediskey.RedisKeyRecordCallFromLockTo, toUid), strconv.FormatInt(fromUid, 10), int(sec))
}

const delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`

func (r *ChatRedis) ReleaseLockTargetUser(ctx context.Context, fromUid, toUid int64) (bool, error) {
	resp, err := r.redisClient.EvalCtx(ctx, delCommand, []string{fmt.Sprintf(rediskey.RedisKeyRecordCallFromLockTo, toUid)}, []string{strconv.FormatInt(fromUid, 10)})
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}
