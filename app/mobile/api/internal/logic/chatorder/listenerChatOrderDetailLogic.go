package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerChatOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListenerChatOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerChatOrderDetailLogic {
	return &ListenerChatOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListenerChatOrderDetailLogic) ListenerChatOrderDetail(req *types.GetListenerChatOrderDetailReq) (resp *types.GetListenerChatOrderDetailResp, err error) {
	var in pb.GetChatOrderDetailReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderDetail(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	var val types.ListenerSeeChatOrder
	_ = copier.Copy(&val, rsp.Order)
	resp = &types.GetListenerChatOrderDetailResp{Info: &val}
	return
}
