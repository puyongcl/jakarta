package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	Cache cache.CacheConf
	Redis redis.RedisConf

	DB struct {
		AdminDataSource string
	}

	ListenerRpcConf   zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	OrderRpcConf      zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	StatRpcConf       zrpc.RpcClientConf
	ImRpcConf         zrpc.RpcClientConf
	PaymentRpcConf    zrpc.RpcClientConf
}
