package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/listener/mq/internal/config"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/redisModel/listenerRedisModel"
	"jakarta/app/redisModel/userRedisModel"
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

	RedisClient   *redis.Redis
	ListenerRedis *listenerRedisModel.ListenerRedis
	UserRedis     *userRedisModel.UserRedis
	ImRedis       *imRedisModel.IMRedis

	AsynqClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()

	return &ServiceContext{
		Config:      c,
		RedisClient: rc,

		OrderRpc:                  order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc:             usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		ChatRpc:                   chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		ListenerRpc:               listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		PaymentRpc:                payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		KqueueSendDefineMsgClient: kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),

		ListenerRedis: listenerRedisModel.NewListenerRedis(rc),
		UserRedis:     userRedisModel.NewUserRedis(rc),
		ImRedis:       imRedisModel.NewIMRedis(rc),

		AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAsynq.Host, Password: c.RedisAsynq.Pass}),
	}
}
