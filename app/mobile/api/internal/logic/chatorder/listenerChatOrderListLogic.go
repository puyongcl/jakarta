package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerChatOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListenerChatOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerChatOrderListLogic {
	return &ListenerChatOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListenerChatOrderListLogic) ListenerChatOrderList(req *types.GetListenerSeeChatOrderListReq) (resp *types.GetListenerSeeChatOrderListResp, err error) {
	var in pb.GetChatOrderListReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerSeeChatOrderListResp{List: make([]*types.ListenerSeeChatOrder, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
