package logic

import (
	"context"
	"fmt"
	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"
	paymentPgModel2 "jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/paykey"
	"jakarta/common/third_party/wxpay"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	//1、payment record confirm
	pf, err := l.svcCtx.ThirdPaymentFlowModel.FindOne(l.ctx, in.FlowNo)
	if err != nil && err != paymentPgModel2.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("UpdateTradeState FindOne db err , sn : %s , err : %+v", in.FlowNo, err))
	}

	if pf == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateTradeState flow no exists, payment flow no: %s", in.FlowNo))
	}

	// Judgment status
	payStatus := wxpay.GetPayStatusByWXPayTradeState(in.TradeState)
	switch payStatus {
	case paykey.ThirdPaymentPayTradeStateUserPaying, paykey.ThirdPaymentPayTradeStateFail, paykey.ThirdPaymentPayTradeStateNotPay, paykey.ThirdPaymentPayTradeStateClosed, paykey.ThirdPaymentPayTradeStateRevoked, paykey.ThirdPaymentPayTradeStatePayError, paykey.ThirdPaymentPayTradeStateSuccess:
		if !tool.IsInt64ArrayExist(pf.PayStatus, paykey.CanUpdatePayResultState) {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateTradeState wrong paystatus:%d", payStatus))
		}

	case paykey.ThirdPaymentPayTradeStateRefundSuccess:
		if pf.PayStatus != paykey.ThirdPaymentPayTradeStateSuccess {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateTradeState wrong paystatus:%d", payStatus))
		}

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateTradeState wrong paystatus:%d", payStatus))
	}

	//3、update .
	newData := new(paymentPgModel2.ThirdPaymentFlow)
	newData.FlowNo = in.FlowNo
	newData.PayStatus = payStatus
	switch payStatus {
	case paykey.ThirdPaymentPayTradeStateSuccess:
		newData.TradeState = in.TradeState
		newData.TransactionId = in.TransactionId
		newData.TradeType = in.TradeType
		newData.TradeStateDesc = in.TradeStateDesc
		newData.BankType = in.BankType
		newData.ActualPayAmount = in.PayAmount
		newData.PayTime.Time = time.Unix(in.PayTime, 0)
		newData.PayTime.Valid = true

		// 比对金额
		if pf.PayAmount != in.PayAmount {
			newData.PayStatus = paykey.ThirdPaymentPayTradeStateSuccessAmountError
		}
	default:
		newData.TradeState = in.TradeState
		newData.TradeStateDesc = in.TradeStateDesc
		newData.PayTime.Time = time.Unix(in.PayTime, 0)
		newData.PayTime.Valid = true
	}

	err = l.svcCtx.ThirdPaymentFlowModel.UpdatePart(l.ctx, newData)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf(" UpdateTradeState db err:%v ,pf : %+v , in : %+v", err, pf, in))
	}
	return &pb.UpdateTradeStateResp{PayStatus: payStatus, OrderId: pf.OrderId, OrderType: pf.OrderType}, nil
}
