package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/common/config"
)

type Config struct {
	service.ServiceConf

	// kq : pub sub
	KqUpdateOrderActionConf kq.KqConf

	// kq client
	KqSendDefineMsgClientConf config.KqConfig
	UploadUserEventClientConf config.KqConfig

	// rpc
	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
	ChatRpcConf       zrpc.RpcClientConf
	ListenerRpcConf   zrpc.RpcClientConf
	PaymentRpcConf    zrpc.RpcClientConf
}
