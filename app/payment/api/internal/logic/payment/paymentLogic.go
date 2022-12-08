package payment

import (
	"context"
	"fmt"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/api/internal/svc"
	"jakarta/app/payment/api/internal/types"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/ctxdata"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/paykey"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentLogic {
	return &PaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentLogic) Payment(req *types.ThirdPaymentWxPayReq) (resp *types.ThirdPaymentWxPayResp, err error) {
	var payAmount int64    // Total amount paid for current order(cent)
	var description string // Current Payment Description.

	payAmount, description, err = l.getPayChatPriceDescription(req.OrderId)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("getPayChatPriceDescription err : %+v req: %+v", err, req))
	}

	// 1、get user openId
	userId := ctxdata.GetUidFromCtx(l.ctx)
	userResp, err := l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		Uid: userId,
	})
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("err:%+v , userId: %d , orderId:%s", err, userId, req.OrderId))
	}
	if userResp.UserAuth == nil || userResp.UserAuth.Uid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 userId: %d , orderId:%s", userId, req.OrderId))
	}

	// 2、create local third payment record
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		Uid:         userId,
		PayModel:    paykey.ThirdPaymentPayModelWechatPay,
		PayAmount:   payAmount,
		OrderId:     req.OrderId,
		OrderType:   req.OrderType,
		Description: description,
		OpenId:      userResp.UserAuth.AuthKey,
	})
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("create local third payment record fail : err: %v , userId: %d,payAmount: %d , orderId: %s", err, userId, payAmount, req.OrderId))
	}

	return &types.ThirdPaymentWxPayResp{
		Appid:     createPaymentResp.Appid,
		NonceStr:  createPaymentResp.NonceStr,
		PaySign:   createPaymentResp.PaySign,
		Package:   createPaymentResp.Package,
		Timestamp: createPaymentResp.Timestamp,
		SignType:  createPaymentResp.SignType,
	}, nil
}

// Get the price and description information of the current order of the paid B&B
func (l *PaymentLogic) getPayChatPriceDescription(orderId string) (int64, string, error) {
	description := "XX服务"
	// get user openid
	resp, err := l.svcCtx.OrderRpc.GetChatOrderDetail(l.ctx, &order.GetChatOrderDetailReq{
		OrderId: orderId,
	})
	if err != nil {
		return 0, description, xerr.NewGrpcErrCodeMsg(xerr.PaymentFail, fmt.Sprintf("OrderRpc.ChatOrderDetail err: %v, orderId: %s", err, orderId))
	}
	description = listenerkey.GetChatOrderDescription(resp.Order.OrderType, resp.Order.BuyUnit*resp.Order.ChatUnitMinute)

	return resp.Order.ActualAmount, description, nil
}
