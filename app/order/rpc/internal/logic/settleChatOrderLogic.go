package logic

import (
	"context"
	"fmt"
	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/money"
	"jakarta/common/tool"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SettleChatOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSettleChatOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettleChatOrderLogic {
	return &SettleChatOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  结算订单
func (l *SettleChatOrderLogic) SettleChatOrder(in *pb.SettleChatOrderReq) (*pb.SettleChatOrderResp, error) {
	// 订单数据
	orderData, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	// 价格分成数据
	pp, err := l.svcCtx.ChatOrderPricingPlanModel.FindOne(l.ctx, orderData.PricingPlanId)
	if err != nil {
		return nil, err
	}
	//
	var confirmUnit int64
	confirmUnit = tool.DivideInt64(orderData.BuyMinuteSum, orderData.ChatUnitMinute)

	// 计算分成
	platformShareAmount, amount := money.CalcShareAmount(in.Star, orderData.BuyUnit, confirmUnit, &money.CurrentShareArg{
		TaxAmount:            orderData.TaxAmount,
		ActualAmount:         orderData.ActualAmount,
		ShareRateStep1Star5:  pp.ShareRateStep1Star5,
		ShareRateStep1Star3:  pp.ShareRateStep1Star3,
		ShareRateStep1Star1:  pp.ShareRateStep1Star1,
		ShareAmountStep1Unit: pp.ShareAmountStep1Unit,
		ShareRateStep2Star5:  pp.ShareRateStep2Star5,
		ShareRateStep2Star3:  pp.ShareRateStep2Star3,
		ShareRateStep2Star1:  pp.ShareRateStep2Star1,
	})

	resp := &pb.SettleChatOrderResp{
		ListenerUid:         orderData.ListenerUid,
		Amount:              amount,
		OrderId:             orderData.OrderId,
		OrderType:           orderData.OrderType,
		UsedMinute:          orderData.UsedChatMinute,
		Uid:                 orderData.Uid,
		PlatformShareAmount: platformShareAmount,
	}
	if in.SettleType != listenerkey.ListenerSettleTypeConfirm { // 不是最终确认 不更新 仅试算
		return resp, nil
	}
	// 最终确认 更新订单已经结算状态
	now := time.Now()
	err = l.svcCtx.ChatOrderModel.SettleOrder(l.ctx, orderData.OrderId, platformShareAmount, amount, orderkey.ChatOrderStateSettle21, &now)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SettleChatOrderLogic SettleOrder err:%+v orderId:%s", err, orderData.OrderId)
	}

	createActionLog(l.svcCtx.ChatOrderStatusLogModel, l.ctx, in.OrderId, orderkey.ChatOrderStateSettle21, orderkey.DefaultSystemOperatorUid, fmt.Sprintf("平台：%d XXX：%d", platformShareAmount, amount), err)
	return resp, nil
}
