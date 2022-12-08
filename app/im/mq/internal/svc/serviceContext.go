package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/im/mq/internal/cacheModel"
	"jakarta/app/im/mq/internal/config"
	"jakarta/app/im/rpc/im"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/statistic/rpc/statistic"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/third_party/zhihu"
)

type ServiceContext struct {
	Config      config.Config
	ImRpc       im.Im
	UserRpc     usercenter.Usercenter
	ListenerRpc listener.Listener
	StatRpc     statistic.Statistic
	ChatRpc     chat.Chat

	RedisClient   *redis.Redis
	IMRedis       *imRedisModel.IMRedis
	ImMemoryCache *cacheModel.IMMemoryCache

	KqueueSendDefineMsgClient          *kq.Pusher
	KqueueSendWxMiniProgramMsgClient   *kq.Pusher
	KqueueSendSubscribeNotifyMsgClient *kq.Pusher
	KqueueSendWxFwhMsgClient           *kq.Pusher
	Zrc                                *zhihu.RestApiClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()
	uRpc := usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf))
	return &ServiceContext{
		Config:      c,
		ImRpc:       im.NewIm(zrpc.MustNewClient(c.ImRpcConf)),
		UserRpc:     uRpc,
		ListenerRpc: listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		StatRpc:     statistic.NewStatistic(zrpc.MustNewClient(c.StatRpcConf)),
		ChatRpc:     chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),

		RedisClient:   rc,
		IMRedis:       imRedisModel.NewIMRedis(rc),
		ImMemoryCache: cacheModel.NewIMMemoryCache(uRpc),

		KqueueSendDefineMsgClient:          kq.NewPusher(c.SendDefineMsgClientConf.Brokers, c.SendDefineMsgClientConf.Topic),
		KqueueSendWxMiniProgramMsgClient:   kq.NewPusher(c.SendWxMiniProgramMsgClientConf.Brokers, c.SendWxMiniProgramMsgClientConf.Topic),
		KqueueSendSubscribeNotifyMsgClient: kq.NewPusher(c.SendSubscribeNotifyMsgClientConf.Brokers, c.SendSubscribeNotifyMsgClientConf.Topic),
		KqueueSendWxFwhMsgClient:           kq.NewPusher(c.SendWxFwhProgramMsgClientConf.Brokers, c.SendWxFwhProgramMsgClientConf.Topic),
		Zrc:                                zhihu.InitZhihuApiClient(),
	}
}
