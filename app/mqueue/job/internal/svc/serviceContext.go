package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/mqueue/job/internal/config"
	"jakarta/app/mqueue/job/internal/redisModel"
	"jakarta/app/order/rpc/order"
	"jakarta/app/statistic/rpc/statistic"
	"jakarta/app/usercenter/rpc/usercenter"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server

	RedisClient *redis.Redis

	JobRedis *redisModel.JobRedis

	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
	ChatRpc       chat.Chat
	ListenerRpc   listener.Listener
	StatRpc       statistic.Statistic

	KqueueSendDefineMsgClient  *kq.Pusher
	KqueueCheckChatStateClient *kq.Pusher
	KqueueUpdateChatStatClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		Config:                     c,
		RedisClient:                rc,
		JobRedis:                   redisModel.NewJobRedis(rc),
		AsynqServer:                newAsynqServer(c),
		OrderRpc:                   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc:              usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		ChatRpc:                    chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		ListenerRpc:                listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		StatRpc:                    statistic.NewStatistic(zrpc.MustNewClient(c.StatRpcConf)),
		KqueueSendDefineMsgClient:  kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		KqueueCheckChatStateClient: kq.NewPusher(c.KqCheckChatStateConf.Brokers, c.KqCheckChatStateConf.Topic),
		KqueueUpdateChatStatClient: kq.NewPusher(c.KqUpdateChatStatConf.Brokers, c.KqUpdateChatStatConf.Topic),
	}
}
