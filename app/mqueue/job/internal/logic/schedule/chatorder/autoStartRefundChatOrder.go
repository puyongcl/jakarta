package chatorder

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"time"
)

type AutoStartRefundChatOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoStartRefundChatOrderHandler(svcCtx *svc.ServiceContext) *AutoStartRefundChatOrderHandler {
	return &AutoStartRefundChatOrderHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoStartRefundChatOrderHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	bTime := time.Now().AddDate(0, 0, -orderkey.AutoStartRefundDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbOrder.GetAutoProcessOrderResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.OrderRpc.GetAutoProcessOrder(ctx, &pbOrder.GetAutoProcessOrderReq{
			PageNo:     pageNo,
			PageSize:   10,
			BeforeTime: bTime,
			State:      orderkey.CanStartRefundOrderState,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoStartRefundChatOrderHandler GetAutoProcessOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoStartRefundChatOrderHandler exit. process order cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoStartRefundChatOrderHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("AutoStartRefundChatOrderHandler done")
	return nil
}

func (l *AutoStartRefundChatOrderHandler) update(ctx context.Context, inList []*pbOrder.AutoProcessOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &pbOrder.DoChatOrderActionReq{
			OrderId:     inList[idx].OrderId,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
			Action:      orderkey.ChatOrderStateAutoStartRefund25,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
