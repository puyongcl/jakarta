package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChatOrderStateLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserChatOrderStateLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChatOrderStateLogLogic {
	return &UserChatOrderStateLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserChatOrderStateLogLogic) UserChatOrderStateLog(req *types.GetChatOrderStateLogReq) (resp *types.GetChatOrderStateLogResp, err error) {
	var in pb.GetChatOrderStateLogReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderStateLog(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetChatOrderStateLogResp{List: make([]*types.ChatOrderStateLog, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
