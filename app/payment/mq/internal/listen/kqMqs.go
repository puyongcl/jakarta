package listen

import (
	"context"
	"jakarta/app/payment/mq/internal/config"
	kqMq "jakarta/app/payment/mq/internal/mqs/kq"
	"jakarta/app/payment/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

//pub sub use kq (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		// 更新支付状态
		kq.MustNewQueue(c.KqUpdatePaymentStatusConf, kqMq.NewUpdatePaymentStatusMq(ctx, svcContext)),
		// 更新退款状态
		kq.MustNewQueue(c.KqUpdateRefundStatusConf, kqMq.NewUpdateRefundStatusMq(ctx, svcContext)),
		// 更新转账状态
		kq.MustNewQueue(c.KqUpdateCashStatusConf, kqMq.NewUpdateCashStatusMq(ctx, svcContext)),
		// 提交转账
		kq.MustNewQueue(c.KqCommitMoveCashConf, kqMq.NewCommitMoveCashMq(ctx, svcContext)),
	}
}
