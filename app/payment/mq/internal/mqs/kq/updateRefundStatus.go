package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/mq/internal/svc"
	"jakarta/app/payment/rpc/pb"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/paykey"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRefundStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRefundStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRefundStatusMq {
	return &UpdateRefundStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRefundStatusMq) Consume(_, val string) error {
	var message kqueue.UpdateRefundStatusMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateRefundStatusMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateRefundStatusMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *UpdateRefundStatusMq) execService(message *kqueue.UpdateRefundStatusMessage) error {
	// Update the flow status.
	req := pb.UpdateRefundReq{
		FlowNo:          message.FlowNo,
		PayFlowNo:       message.PayFlowNo,
		TransactionId:   message.TransactionId,
		ReceivedAccount: message.ReceivedAccount,
		SuccessTime:     message.SuccessTime,
		Status:          message.Status,
		Amount:          message.Amount,
	}
	rsp, err := l.svcCtx.PaymentRpc.UpdateRefundState(l.ctx, &req)
	if err != nil {
		return err
	}
	orderTradeState := l.getOrderStateByRefundState(rsp.RefundStatus)
	switch orderTradeState {
	case orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateFinishRefund8, orderkey.ChatOrderStatePayFail22, orderkey.ChatOrderStateRefundPayFail23:
		break
	default:
		// 暂不需要处理的状态
		return nil
	}

	// 更新订单支付状态
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &order.DoChatOrderActionReq{
		OrderId: rsp.OrderId,
		Action:  orderTradeState,
	})
	if err != nil {
		return xerr.NewGrpcErrCodeMsg(xerr.ServerCommonError, fmt.Sprintf("UpdateRefundStatusMq ChatOrderAction err: %+v ,message:%+v", err, message))
	}
	return nil
}

//Get order status based on payment status.
func (l *UpdateRefundStatusMq) getOrderStateByRefundState(status int64) int64 {
	switch status {
	case paykey.ThirdPaymentPayTradeStateSuccess:
		return orderkey.ChatOrderStatePaySuccess3
	case paykey.ThirdPaymentPayTradeStateFail:
		return orderkey.ChatOrderStatePayFail22
	case paykey.ThirdPaymentPayTradeStateRefundSuccess:
		return orderkey.ChatOrderStateFinishRefund8
	case paykey.ThirdPaymentPayTradeStateRefundFail:
		return orderkey.ChatOrderStateRefundPayFail23
	default:
		return 0
	}
}
