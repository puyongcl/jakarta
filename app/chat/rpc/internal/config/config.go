package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	DBRO struct {
		DataSource string
	}
	Cache cache.CacheConf

	KqUpdateChatStatConf config.KqConfig
	KqSendDefineMsgConf  config.KqConfig
	KqFirstEnterChatConf config.KqConfig
}
