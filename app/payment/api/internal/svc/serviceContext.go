package svc

import (
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/api/internal/config"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/usercenter/rpc/usercenter"
)

type ServiceContext struct {
	Config config.Config

	KqueueUpdatePaymentStatusClient *kq.Pusher
	KqueueUpdateRefundStatusClient  *kq.Pusher
	KqueueUpdateCashStatusClient    *kq.Pusher

	PaymentRpc    payment.Payment
	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter

	WxpayNotifyHandler *notify.Handler
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		KqueueUpdatePaymentStatusClient: kq.NewPusher(c.KqUpdatePaymentStatusConf.Brokers, c.KqUpdatePaymentStatusConf.Topic),
		KqueueUpdateRefundStatusClient:  kq.NewPusher(c.KqUpdateRefundStatusConf.Brokers, c.KqUpdateRefundStatusConf.Topic),
		KqueueUpdateCashStatusClient:    kq.NewPusher(c.KqUpdateCashStatusConf.Brokers, c.KqUpdateCashStatusConf.Topic),

		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),

		WxpayNotifyHandler: NewWxpayNotifyHandler(&c.WxPayCallbackConf),
	}
}
