package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EstimateChatOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEstimateChatOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimateChatOrderLogic {
	return &EstimateChatOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EstimateChatOrderLogic) EstimateChatOrder(req *types.EstimateOrderReq) (resp *types.EstimateOrderResp, err error) {
	var in pb.EstimateOrderReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.EstimateOrder(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.EstimateOrderResp{}
	_ = copier.Copy(resp, rsp)
	return
}
