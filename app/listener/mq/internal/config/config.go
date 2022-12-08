package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	service.ServiceConf

	Redis      redis.RedisConf
	RedisAsynq redis.RedisConf
	// kq : pub sub
	KqUpdateListenerUserStatConf kq.KqConf
	KqSendHelloWhenUserLoginConf kq.KqConf

	// kq client
	KqSendDefineMsgConf config.KqConfig

	// rpc
	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	ListenerRpcConf   zrpc.RpcClientConf
	PaymentRpcConf    zrpc.RpcClientConf
}
