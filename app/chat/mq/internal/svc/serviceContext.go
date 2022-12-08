package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/chat/mq/internal/config"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/im/rpc/im"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/order/rpc/order"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/usercenter/rpc/usercenter"
)

type ServiceContext struct {
	Config        config.Config
	ImRpc         im.Im
	UsercenterRpc usercenter.Usercenter
	ChatRpc       chat.Chat
	OrderRpc      order.Order
	ListenerRpc   listener.Listener

	RedisClient *redis.Redis
	IMRedis     *imRedisModel.IMRedis

	KqueueSendDefineMsgClient *kq.Pusher
	AsynqClient               *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		Config:        c,
		ImRpc:         im.NewIm(zrpc.MustNewClient(c.ImRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		ChatRpc:       chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		ListenerRpc:   listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),

		RedisClient: rc,
		IMRedis:     imRedisModel.NewIMRedis(rc),

		KqueueSendDefineMsgClient: kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		AsynqClient:               asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAsynq.Host, Password: c.RedisAsynq.Pass}),
	}
}
