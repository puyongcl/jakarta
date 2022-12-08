package third

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"jakarta/app/payment/api/internal/svc"
	"jakarta/app/payment/api/internal/types"
	"jakarta/common/key/db"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type WxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxPayCallbackLogic {
	return &WxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxPayCallbackLogic) WxPayCallback(w http.ResponseWriter, r *http.Request) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	//Verifying signatures, parsing data
	transaction := new(payments.Transaction)
	_, err = l.svcCtx.WxpayNotifyHandler.ParseNotifyRequest(l.ctx, r, transaction)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("WxPayCallbackLogic ParseNotifyRequest err:%+v", err))
	}

	logx.WithContext(l.ctx).Infof("WxRefundCallback req:%+v", transaction)

	returnCode := "SUCCESS"
	msg := ""

	// Update the flow status.
	tt, err := time.Parse(time.RFC3339, *transaction.SuccessTime)
	if err != nil {
		tt = time.Now()
		logx.WithContext(l.ctx).Errorf("wrong pay time layout tran:%+v", transaction)
	}
	kqMsg := kqueue.UpdatePaymentStatusMessage{
		FlowNo:         *transaction.OutTradeNo,
		TradeState:     *transaction.TradeState,
		TransactionId:  *transaction.TransactionId,
		TradeType:      *transaction.TradeType,
		TradeStateDesc: *transaction.TradeStateDesc,
		BankType:       *transaction.BankType,
		PayTime:        tt.Unix(),
		PayAmount:      *transaction.Amount.PayerTotal,
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
	logx.WithContext(l.ctx).Infof("WxRefundCallback push done flow:%s time:%s", kqMsg.FlowNo, time.Now().Format(db.DateTimeFormat))
	return resp, err
}

func (l *WxPayCallbackLogic) pushKq(msg *kqueue.UpdatePaymentStatusMessage) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return l.svcCtx.KqueueUpdatePaymentStatusClient.Push(string(buf))
}
