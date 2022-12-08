package third

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/payment/api/internal/svc"
	"jakarta/app/payment/api/internal/types"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
	"net/http"
)

type WxRefundCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxRefundCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxRefundCallbackLogic {
	return &WxRefundCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxRefundCallbackLogic) WxRefundCallback(w http.ResponseWriter, r *http.Request) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	//Verifying signatures, parsing data
	rf := new(Refund)
	_, err = l.svcCtx.WxpayNotifyHandler.ParseNotifyRequest(l.ctx, r, rf)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("WxRefundCallbackLogic ParseNotifyRequest err:%+v", err))
	}

	returnCode := "SUCCESS"
	msg := ""

	// Update the flow status.
	//logx.WithContext(l.ctx).Infof("WxRefundCallback %+v", *rf)
	kqMsg := kqueue.UpdateRefundStatusMessage{
		FlowNo:          *rf.OutRefundNo,
		PayFlowNo:       *rf.OutTradeNo,
		TransactionId:   *rf.TransactionId,
		ReceivedAccount: *rf.UserReceivedAccount,
		SuccessTime:     rf.SuccessTime.Unix(),
		Amount:          *rf.Amount.PayerRefund,
	}
	if rf.RefundStatus != nil {
		kqMsg.Status = string(*(rf.RefundStatus))
	}
	err = l.pushKq(&kqMsg)
	if err != nil {
		returnCode = "FAIL"
		msg = fmt.Sprintf("%+v", err)
	}

	resp = &types.ThirdPaymentWxPayCallbackResp{
		Code:    returnCode,
		Message: msg,
	}
	return resp, err
}

func (l *WxRefundCallbackLogic) pushKq(msg *kqueue.UpdateRefundStatusMessage) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return l.svcCtx.KqueueUpdateRefundStatusClient.Push(string(buf))
}
