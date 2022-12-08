package kq

import (
	"context"
	"encoding/json"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/mq/internal/svc"
	"jakarta/app/payment/rpc/payment"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/paykey"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaymentStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentStatusMq {
	return &UpdatePaymentStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaymentStatusMq) Consume(_, val string) error {
	var message kqueue.UpdatePaymentStatusMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdatePaymentStatusMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdatePaymentStatusMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *UpdatePaymentStatusMq) execService(message *kqueue.UpdatePaymentStatusMessage) error {
	logx.WithContext(l.ctx).Infof("UpdatePaymentStatusMq time:%s req:%+v", time.Now().Format(db.DateTimeFormat), message)
	rsp, err := l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
		FlowNo:         message.FlowNo,
		TradeState:     message.TradeState,
		TransactionId:  message.TransactionId,
		TradeType:      message.TradeType,
		TradeStateDesc: message.TradeStateDesc,
		BankType:       message.BankType,
		PayTime:        message.PayTime,
		PayAmount:      message.PayAmount,
	})
	if err != nil {
		return err
	}

	orderTradeState := l.getOrderStateByPaymentState(rsp.PayStatus, rsp.OrderType)
	if orderTradeState == 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	// 更新订单支付状态
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &order.DoChatOrderActionReq{
		OrderId:     rsp.OrderId,
		Action:      orderTradeState,
		OperatorUid: orderkey.DefaultSystemOperatorUid,
	})
	if err != nil {
		return err
	}
	return nil
}

//Get order status based on payment status.
func (l *UpdatePaymentStatusMq) getOrderStateByPaymentState(thirdPaymentPayStatus int64, orderType int64) int64 {
	switch thirdPaymentPayStatus {
	case paykey.ThirdPaymentPayTradeStateSuccess:
		return orderkey.GetPaySuccessOrderState(orderType)

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
