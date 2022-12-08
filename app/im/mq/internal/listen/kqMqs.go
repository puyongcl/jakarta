package listen

import (
	"context"
	"jakarta/app/im/mq/internal/config"
	kqMq "jakarta/app/im/mq/internal/mqs/kq"
	"jakarta/app/im/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

//pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
	return []service.Service{
		//Listening for changes in consumption flow status
		kq.MustNewQueue(c.IMDefineMsgSendConf, kqMq.NewSendDefineImMsgMq(ctx, svcContext)),
		kq.MustNewQueue(c.IMStateChangeMsgConf, kqMq.NewImCallbackMq(ctx, svcContext)),
		kq.MustNewQueue(c.SubscribeNotifyMsgConf, kqMq.NewSubscribeNotifyMsgMq(ctx, svcContext)),
		kq.MustNewQueue(c.IMAfterSendMsgConf, kqMq.NewImAfterSendMsgMq(ctx, svcContext)),
		kq.MustNewQueue(c.WxFwhCallbackEventConf, kqMq.NewWxFwhCallbackEventMq(ctx, svcContext)),
		kq.MustNewQueue(c.WxMiniProgramMsgSendConf, kqMq.NewSendMiniProgramSubscribeMsgMq(ctx, svcContext)),
		kq.MustNewQueue(c.WxFwhMsgSendConf, kqMq.NewSendWxFwhNotifyMsgMq(ctx, svcContext)),
		kq.MustNewQueue(c.UploadUserEventConf, kqMq.NewUploadUserEventMq(ctx, svcContext)),
		//.....
	}
}
