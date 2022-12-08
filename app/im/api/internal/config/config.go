package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	IMRpcConf zrpc.RpcClientConf

	KqImStateChangeMsgConf config.KqConfig
	KqImAfterSendMsgConf   config.KqConfig
	WxFwhCallbackEventConf config.KqConfig

	RedisCache redis.RedisConf

	WxFwhConf  config.WxFwhConf
	WxMiniConf config.WxMiniConf
}
