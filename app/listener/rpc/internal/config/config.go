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
	Cache cache.CacheConf

	// kq client
	KqSendDefineMsgConf          config.KqConfig
	KqSubscribeNotifyMsgConf     config.KqConfig
	KqUpdateListenerUserStatConf config.KqConfig
	KqSendHelloWhenUserLoginConf config.KqConfig

	TimConf config.TimConf

	RedisAsynq redis.RedisConf

	TencentConf config.CosConf

	HfbfCashConf config.HfbfCashConfig

	ContractConfig config.ContractConfig
}
