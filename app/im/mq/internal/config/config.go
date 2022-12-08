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
	IMDefineMsgSendConf      kq.KqConf
	IMStateChangeMsgConf     kq.KqConf
	SubscribeNotifyMsgConf   kq.KqConf
	IMAfterSendMsgConf       kq.KqConf
	WxFwhCallbackEventConf   kq.KqConf
	WxMiniProgramMsgSendConf kq.KqConf
	WxFwhMsgSendConf         kq.KqConf
	UploadUserEventConf      kq.KqConf

	// kq client
	SendDefineMsgClientConf          config.KqConfig
	SendWxMiniProgramMsgClientConf   config.KqConfig
	SendWxFwhProgramMsgClientConf    config.KqConfig
	SendSubscribeNotifyMsgClientConf config.KqConfig
	// rpc
	ImRpcConf         zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	ListenerRpcConf   zrpc.RpcClientConf
	StatRpcConf       zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
}
