package config

import (
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

	WxPayCallbackConf config.WxPayCallbackConf

	KqUpdatePaymentStatusConf config.KqConfig
	KqUpdateRefundStatusConf  config.KqConfig
	KqUpdateCashStatusConf    config.KqConfig

	PaymentRpcConf    zrpc.RpcClientConf
	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf

	HfbfCashConf config.HfbfCashConfig
}
