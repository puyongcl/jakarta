package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChatOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChatOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChatOrderDetailLogic {
	return &UserChatOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChatOrderDetailLogic) UserChatOrderDetail(req *types.GetUserChatOrderDetailReq) (resp *types.GetUserChatOrderDetailResp, err error) {
	var in pb.GetChatOrderDetailReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderDetail(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	var val types.UserSeeChatOrder
	_ = copier.Copy(&val, rsp.Order)
	resp = &types.GetUserChatOrderDetailResp{Info: &val}
	return
}
