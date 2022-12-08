package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/paykey"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRequestRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestRefundLogic {
	return &RequestRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发起退款
func (l *RequestRefundLogic) RequestRefund(in *pb.RequestRefundReq) (*pb.RequestRefundResp, error) {
	// 查询当前订单是否有流水
	cnt, err := l.svcCtx.ThirdRefundFlowModel.FindCount(l.ctx, in.OrderId, paykey.NotNeedRefundState)
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentErrorFlowAlreadyExist, "已经退款")
	}

	pf, err := l.svcCtx.ThirdPaymentFlowModel.FindOneByOrderId(l.ctx, in.OrderId, paykey.CanUpdateRefundState)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("%+v", err))
	}
	//
	data := &paymentPgModel.ThirdRefundFlow{
		FlowNo:             uniqueid.GenSn(uniqueid.SnPrefixThirdRefundFlowNo),
		PayFlowNo:          pf.FlowNo,
		TransactionId:      pf.TransactionId,
		OrderId:            pf.OrderId,
		Reason:             in.Reason,
		PayAmount:          pf.ActualPayAmount,
		RefundAmount:       pf.ActualPayAmount,
		ActualRefundAmount: 0,
		RefundStatus:       paykey.ThirdPaymentPayTradeStateStartRefund,
		RefundTime:         sql.NullTime{},
		Uid:                pf.Uid,
	}

	// 退款
	rsvc := refunddomestic.RefundsApiService{Client: l.svcCtx.WxPayClient}
	req := refunddomestic.CreateRequest{
		TransactionId: core.String(data.TransactionId),
		OutTradeNo:    core.String(data.PayFlowNo),
		OutRefundNo:   core.String(data.FlowNo),
		Reason:        core.String(data.Reason),
		NotifyUrl:     core.String(l.svcCtx.Config.WxPayConf.RefundNotifyUrl),
		Amount: &refunddomestic.AmountReq{
			Currency: core.String("CNY"),
			Refund:   core.Int64(data.RefundAmount),
			Total:    core.Int64(data.PayAmount),
		},
	}

	_, _, err = rsvc.Create(l.ctx, req)
	if err != nil {
		data.Remark = fmt.Sprintf("%+v", err)
		logx.WithContext(l.ctx).Errorf("RequestRefundLogic RequestRefund rsvc.Create req:%+v,err:%+v", req, err)
		err = nil
	}

	_, err = l.svcCtx.ThirdRefundFlowModel.Insert(l.ctx, data)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RequestRefundLogic RequestRefund Insert data:%+v,err:%+v", data, err)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("%+v", err))
	}

	return &pb.RequestRefundResp{}, nil
}
