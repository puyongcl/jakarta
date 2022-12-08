package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DBRO struct {
		UserDataSource     string
		OrderDataSource    string
		ListenerDataSource string
	}
	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
}
