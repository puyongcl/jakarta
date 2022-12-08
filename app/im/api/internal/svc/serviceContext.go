package svc

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/im/api/internal/config"
	"jakarta/app/im/rpc/im"
)

type ServiceContext struct {
	Config config.Config
	ImRpc  im.Im

	KqueueImStateChangeMsgClient   *kq.Pusher
	KqueueImAfterSendMsgClient     *kq.Pusher
	KqueueWxFwhCallbackEventClient *kq.Pusher
	Wx                             *wechat.Wechat
	Wxfwh                          *officialaccount.OfficialAccount
}

func NewServiceContext(c config.Config) *ServiceContext {
	wc, fwh := InitWechat(&c)
	return &ServiceContext{
		Config: c,
		ImRpc:  im.NewIm(zrpc.MustNewClient(c.IMRpcConf)),

		KqueueImStateChangeMsgClient:   kq.NewPusher(c.KqImStateChangeMsgConf.Brokers, c.KqImStateChangeMsgConf.Topic),
		KqueueImAfterSendMsgClient:     kq.NewPusher(c.KqImAfterSendMsgConf.Brokers, c.KqImAfterSendMsgConf.Topic),
		KqueueWxFwhCallbackEventClient: kq.NewPusher(c.WxFwhCallbackEventConf.Brokers, c.WxFwhCallbackEventConf.Topic),
		Wx:                             wc,
		Wxfwh:                          fwh,
	}
}
