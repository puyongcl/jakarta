package chatRedisModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"testing"
)

func TestListenerRedisModel_GetSetListenerFreeMsgCnt(t *testing.T) {
	rcf := redis.RedisConf{
		Host: "192.168.1.12:36379",
		Type: "node",
		Pass: "G62m50oigInC30sf",
	}
	rc := rcf.NewRedis()
	err := rc.Set("test", "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	rn, err := rc.Decr("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rn)

	rn, err = rc.Decr("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rn)

	rs, err := rc.Get("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rs)
}

func TestListenerRedisModel_ZSore(t *testing.T) {
	rcf := redis.RedisConf{
		Host: "192.168.1.94:36379",
		Type: "node",
		Pass: "G62m50oigInC30sf",
	}
	rc := rcf.NewRedis()
	ctx := context.Background()
	var err error
	_, err = rc.ZaddCtx(ctx, "test1", 0, "t1")
	if err != nil {
		fmt.Println(err)
	}
	_, err = rc.ZaddCtx(ctx, "test1", 0, "t2")
	if err != nil {
		fmt.Println(err)
	}
	_, err = rc.ZaddCtx(ctx, "test1", 0, "t3")
	if err != nil {
		fmt.Println(err)
	}
	_, err = rc.ZaddCtx(ctx, "test1", 0, "t4")
	if err != nil {
		fmt.Println(err)
	}
	var v int64
	v, err = rc.ZrankCtx(ctx, "test1", "t4")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
}
