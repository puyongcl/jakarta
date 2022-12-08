package listen

import (
	"context"
	"jakarta/app/listener/mq/internal/config"
	kqMq "jakarta/app/listener/mq/internal/mqs/kq"
	"jakarta/app/listener/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

//pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		// 更新XXX和用户的交互情况
		kq.MustNewQueue(c.KqUpdateListenerUserStatConf, kqMq.NewUpdateListenerUserStatMq(ctx, svcContext)),
		// 用户登陆推荐XXX
		kq.MustNewQueue(c.KqSendHelloWhenUserLoginConf, kqMq.NewSendHelloWhenUserLoginMq(ctx, svcContext)),
	}
}
