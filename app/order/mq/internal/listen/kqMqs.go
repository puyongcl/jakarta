package listen

import (
	"context"
	"jakarta/app/order/mq/internal/config"
	kqMq "jakarta/app/order/mq/internal/mqs/kq"
	"jakarta/app/order/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

//pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		// 订单状态变更的后续操作
		kq.MustNewQueue(c.KqUpdateOrderActionConf, kqMq.NewUpdateOrderActionMq(ctx, svcContext)),
	}
}
