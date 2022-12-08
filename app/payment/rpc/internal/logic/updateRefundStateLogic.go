package logic

import (
	"context"
	"fmt"
	paymentPgModel2 "jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/paykey"
	"jakarta/common/third_party/wxpay"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRefundStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRefundStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRefundStateLogic {
	return &UpdateRefundStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新退款状态
func (l *UpdateRefundStateLogic) UpdateRefundState(in *pb.UpdateRefundReq) (*pb.UpdateRefundResp, error) {
	//1、payment record confirm
	pf, err := l.svcCtx.ThirdRefundFlowModel.FindOne(l.ctx, in.FlowNo)
	if err != nil && err != paymentPgModel2.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("UpdateRefundState FindOneBySn db err , sn : %s , err : %+v", in.FlowNo, err))
	}

	if pf == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateRefundState not find payment flow no: %s", in.FlowNo))
	}

	// Judgment status
	refundStatus := wxpay.GetRefundStatusByWXPayTradeState(in.Status)
	switch refundStatus {
	case paykey.ThirdPaymentPayTradeStateRefundSuccess, paykey.ThirdPaymentPayTradeStateRefundClosed, paykey.ThirdPaymentPayTradeStateRefundFail:
		if !tool.IsInt64ArrayExist(pf.RefundStatus, paykey.CanUpdateRefundResult) {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateRefundState wrong paystatus:%d", refundStatus))
		}

	case paykey.ThirdPaymentPayTradeStateRefundProcessing:
		if pf.RefundStatus != paykey.ThirdPaymentPayTradeStateStartRefund {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateRefundState wrong paystatus:%d", refundStatus))
		}

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UpdateRefundState wrong paystatus:%d", refundStatus))
	}

	//3、update .
	newData := new(paymentPgModel2.ThirdRefundFlow)
	newData.RefundStatus = refundStatus
	newData.FlowNo = in.FlowNo
	switch refundStatus {
	case paykey.ThirdPaymentPayTradeStateRefundSuccess:
		newData.TransactionId = in.TransactionId
		newData.ActualRefundAmount = in.Amount
		newData.WxStatus = in.Status
		newData.ReceivedAccount = in.ReceivedAccount
		newData.RefundTime.Time = time.Unix(in.SuccessTime, 0)
		newData.RefundTime.Valid = true

		// 比对金额
		if pf.RefundAmount != in.Amount {
			newData.RefundStatus = paykey.ThirdPaymentPayTradeStateRefundSuccessAmountError
		}
	default:
		newData.WxStatus = in.Status
		newData.RefundTime.Time = time.Unix(in.SuccessTime, 0)
		newData.RefundTime.Valid = true
	}

	err = l.svcCtx.ThirdRefundFlowModel.UpdatePart(l.ctx, newData)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("UpdateRefundState db err:%+v ,pf : %+v , in : %+v", err, pf, in))
	}

	return &pb.UpdateRefundResp{RefundStatus: refundStatus, OrderId: pf.OrderId}, nil
}
