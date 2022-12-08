package delay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/app/order/rpc/order"
	"jakarta/common/key/orderkey"
	"jakarta/common/xerr"
)

// CloseChatOrderHandler close no pay chatOrder
type CloseChatOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCloseChatOrderHandler(svcCtx *svc.ServiceContext) *CloseChatOrderHandler {
	return &CloseChatOrderHandler{
		svcCtx: svcCtx,
	}
}

// defer  close no pay chatOrder  : if return err != nil , asynq will retry
func (l *CloseChatOrderHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p jobtype.DeferCloseChatOrderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("closeChatOrderStateMqHandler payload err:%+v, payLoad:%+v", err, t.Payload()))
	}

	resp, err := l.svcCtx.OrderRpc.GetChatOrderDetail(ctx, &order.GetChatOrderDetailReq{
		OrderId: p.OrderId,
	})
	if err != nil {
		return xerr.NewGrpcErrCodeMsg(xerr.ServerCommonError, fmt.Sprintf("closeChatOrderStateMqHandler get order fail or order no exists err:%v, orderId:%s", err, p.OrderId))
	}

	if resp.Order.OrderState == orderkey.ChatOrderStateWaitPay1 {
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &order.DoChatOrderActionReq{
			OrderId:     p.OrderId,
			Action:      orderkey.ChatOrderStateCancel2,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
		})
		if err != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("CloseChatOrderHandler close order fail err:%v, sn:%s ", err, p.OrderId))
		}
	}

	return nil
}
