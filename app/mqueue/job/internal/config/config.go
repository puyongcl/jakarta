package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	service.ServiceConf
	RedisAsynq redis.RedisConf
	Redis      redis.RedisConf

	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	ListenerRpcConf   zrpc.RpcClientConf
	StatRpcConf       zrpc.RpcClientConf

	KqSendDefineMsgConf  config.KqConfig
	KqCheckChatStateConf config.KqConfig
	KqUpdateChatStatConf config.KqConfig
}
