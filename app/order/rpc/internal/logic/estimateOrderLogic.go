package logic

import (
	"context"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EstimateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEstimateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimateOrderLogic {
	return &EstimateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取订单费用估算
func (l *EstimateOrderLogic) EstimateOrder(in *pb.EstimateOrderReq) (*pb.EstimateOrderResp, error) {
	// 定价方案
	var pp *pb.GetBusinessChatPricePlanResp
	var err error
	gpp := NewGetBusinessChatPricePlanLogic(l.ctx, l.svcCtx)
	pp, err = gpp.getListenerPrice(in.Uid, in.PricingPlanId)
	if err != nil {
		return nil, err
	}
	// 估算
	now := time.Now()
	rsp := calculate(in.TextChatPrice, in.VoiceChatPrice, in.OrderType, in.BuyUnit, &now, pp.Config)
	return rsp, nil
}
