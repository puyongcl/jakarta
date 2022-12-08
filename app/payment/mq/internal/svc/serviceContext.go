package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/mq/internal/config"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/usercenter/rpc/usercenter"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
	ChatRpc       chat.Chat
	ListenerRpc   listener.Listener
	PaymentRpc    payment.Payment

	KqueueSendDefineMsgClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		ChatRpc:       chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		ListenerRpc:   listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),

		KqueueSendDefineMsgClient: kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
	}
}
