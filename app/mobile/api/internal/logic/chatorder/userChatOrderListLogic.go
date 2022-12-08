package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChatOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChatOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChatOrderListLogic {
	return &UserChatOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChatOrderListLogic) UserChatOrderList(req *types.GetUserChatOrderListReq) (resp *types.GetUserChatOrderListResp, err error) {
	var in pb.GetChatOrderListReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetUserChatOrderListResp{List: make([]*types.UserSeeChatOrder, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
