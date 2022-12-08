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

	Redis redis.RedisConf

	// kq : pub sub
	UpdateChatStatConf kq.KqConf
	CheckChatStateConf kq.KqConf
	FirstEnterChatConf kq.KqConf

	// rpc
	ImRpcConf         zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	OrderRpcConf      zrpc.RpcClientConf
	ListenerRpcConf   zrpc.RpcClientConf

	KqSendDefineMsgConf config.KqConfig

	RedisAsynq redis.RedisConf
}
