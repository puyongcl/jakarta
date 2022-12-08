package order

import (
	"context"
	"github.com/jinzhu/copier"
	pbOrder "jakarta/app/order/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRefundOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundOrderListLogic {
	return &GetRefundOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRefundOrderListLogic) GetRefundOrderList(req *types.GetRefundOrderListReq) (resp *types.GetRefundOrderListResp, err error) {
	var in pbOrder.GetOrderListByAdminReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.OrderRpc.GetOrderListByAdmin(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetRefundOrderListResp{
		Sum:  rs.Sum,
		List: make([]*types.RefundChatOrder, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
