package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBusinessChatPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBusinessChatPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBusinessChatPriceLogic {
	return &GetBusinessChatPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBusinessChatPriceLogic) GetBusinessChatPrice(req *types.GetBusinessChatPricingPlanReq) (resp *types.GetBusinessChatPricingPlanResp, err error) {
	//
	rsp, err := l.svcCtx.OrderRpc.GetBusinessChatPricePlan(l.ctx, &pb.GetBusinessChatPricePlanReq{Uid: req.Uid})
	if err != nil {
		return nil, err
	}
	var val types.BusinessChatPricingPlan
	_ = copier.Copy(&val, rsp.Config)
	resp = &types.GetBusinessChatPricingPlanResp{Config: &val}
	return
}
