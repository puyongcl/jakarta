package listen

import (
	"context"
	"jakarta/app/chat/mq/internal/config"
	kqMq "jakarta/app/chat/mq/internal/mqs/kq"
	"jakarta/app/chat/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

//pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
	return []service.Service{
		//Listening for changes in consumption flow status
		kq.MustNewQueue(c.UpdateChatStatConf, kqMq.NewUpdateChatStatMq(ctx, svcContext)),
		kq.MustNewQueue(c.CheckChatStateConf, kqMq.NewCheckChatStateMq(ctx, svcContext)),
		kq.MustNewQueue(c.FirstEnterChatConf, kqMq.NewUserEnterChatMq(ctx, svcContext)),
		//.....
	}
}
