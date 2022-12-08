package svc

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/bbs/rpc/bbs"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/mobile/api/internal/config"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/statistic/rpc/statistic"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/third_party/tencentcloud"
)

type ServiceContext struct {
	Config        config.Config
	UsercenterRpc usercenter.Usercenter
	OrderRpc      order.Order
	ListenerRpc   listener.Listener
	ChatRpc       chat.Chat
	PaymentRpc    payment.Payment
	StatRpc       statistic.Statistic
	BbsRpc        bbs.Bbs

	TencentAPIClient *tencentcloud.TencentStsClient
	WxMini           *miniprogram.MiniProgram

	KqueueSendDefineMsgClient  *kq.Pusher
	KqSubscribeNotifyMsgClient *kq.Pusher

	RedisClient *redis.Redis
	IMRedis     *imRedisModel.IMRedis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		ListenerRpc:   listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		ChatRpc:       chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		StatRpc:       statistic.NewStatistic(zrpc.MustNewClient(c.StatRpcConf)),
		BbsRpc:        bbs.NewBbs(zrpc.MustNewClient(c.BbsRpcConf)),

		TencentAPIClient: tencentcloud.InitStsClient(c.TencentConf.SecretId, c.TencentConf.SecretKey, c.TencentConf.AppId, c.TencentConf.Bucket, c.TencentConf.Region, c.TencentConf.Expire),
		WxMini: wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
			AppID:     c.WxMiniConf.AppId,
			AppSecret: c.WxMiniConf.Secret,
			Cache: cache.NewRedis(context.Background(), &cache.RedisOpts{
				Host:     c.RedisCache.Host,
				Password: c.RedisCache.Pass,
			}),
		}),

		KqueueSendDefineMsgClient:  kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		KqSubscribeNotifyMsgClient: kq.NewPusher(c.KqSubscribeNotifyMsgConf.Brokers, c.KqSubscribeNotifyMsgConf.Topic),

		RedisClient: rc,
		IMRedis:     imRedisModel.NewIMRedis(rc),
	}
}
