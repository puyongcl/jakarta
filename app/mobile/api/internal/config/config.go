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
	ListenerRpcConf   zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	PaymentRpcConf    zrpc.RpcClientConf
	StatRpcConf       zrpc.RpcClientConf
	BbsRpcConf        zrpc.RpcClientConf

	AppVerConf  config.AppVerConf
	WxMiniConf  config.WxMiniConf
	TencentConf config.CosConf

	KqSendDefineMsgConf      config.KqConfig
	KqSubscribeNotifyMsgConf config.KqConfig

	RedisCache redis.RedisConf
	Redis      redis.RedisConf
}
