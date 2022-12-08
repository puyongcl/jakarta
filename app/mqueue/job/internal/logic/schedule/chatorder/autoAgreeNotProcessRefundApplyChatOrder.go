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

type AutoAgreeNotProcessRefundApplyChatOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoAgreeNotProcessRefundApplyChatOrderHandler(svcCtx *svc.ServiceContext) *AutoAgreeNotProcessRefundApplyChatOrderHandler {
	return &AutoAgreeNotProcessRefundApplyChatOrderHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoAgreeNotProcessRefundApplyChatOrderHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	bTime := time.Now().AddDate(0, 0, -orderkey.AutoAgreeNotProcessRefundApplyChatOrderDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbOrder.GetAutoProcessOrderResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.OrderRpc.GetAutoProcessOrder(ctx, &pbOrder.GetAutoProcessOrderReq{
			PageNo:     pageNo,
			PageSize:   10,
			BeforeTime: bTime,
			State:      []int64{orderkey.ChatOrderStateApplyRefund5},
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoAgreeNotProcessRefundApplyChatOrderHandler GetAutoProcessOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoAgreeNotProcessRefundApplyChatOrderHandler exit. process order cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoAgreeNotProcessRefundApplyChatOrderHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("AutoAgreeNotProcessRefundApplyChatOrderHandler done")
	return nil
}

func (l *AutoAgreeNotProcessRefundApplyChatOrderHandler) update(ctx context.Context, inList []*pbOrder.AutoProcessOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &pbOrder.DoChatOrderActionReq{
			OrderId:     inList[idx].OrderId,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
			Action:      orderkey.ChatOrderStateAutoAgreeNotProcessRefund27,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
