package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	Cache      cache.CacheConf
	RedisAsynq redis.RedisConf

	KqUpdateOrderActionConf config.KqConfig
	// kq client
	KqSendDefineMsgConf config.KqConfig

	KqSendSubscribeNotifyMsgClientConf config.KqConfig
}
