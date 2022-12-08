package chatorder

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/tool"
	"time"
)

type AutoConfirmFinishChatOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoConfirmFinishChatOrderHandler(svcCtx *svc.ServiceContext) *AutoConfirmFinishChatOrderHandler {
	return &AutoConfirmFinishChatOrderHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoConfirmFinishChatOrderHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	bTime := time.Now().AddDate(0, 0, -orderkey.AutoFinishAfterNotGoodCommentDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbOrder.GetAutoProcessOrderResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.OrderRpc.GetAutoProcessOrder(ctx, &pbOrder.GetAutoProcessOrderReq{
			PageNo:     pageNo,
			PageSize:   10,
			BeforeTime: bTime,
			State:      orderkey.CanAutoConfirmFinishOrderState,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoConfirmFinishChatOrderHandler GetAutoProcessOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoConfirmFinishChatOrderHandler exit. process order cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoConfirmFinishChatOrderHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("AutoConfirmFinishChatOrderHandler done")
	return nil
}

func (l *AutoConfirmFinishChatOrderHandler) update(ctx context.Context, inList []*pbOrder.AutoProcessOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		var in = pbOrder.DoChatOrderActionReq{
			OrderId:     inList[idx].OrderId,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
			Action:      orderkey.ChatOrderStateAutoConfirmFinish19,
		}
		if tool.IsInt64ArrayExist(inList[idx].OrderState, orderkey.RefuseRefundOrderState) { // 非一般评价
			in.Action = orderkey.ChatOrderStateRefundRefuseAutoFinish20
		}
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &in)
		if err != nil {
			return err
		}
	}

	return nil
}
