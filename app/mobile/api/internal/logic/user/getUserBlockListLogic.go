package user

import (
	"context"
	"github.com/jinzhu/copier"
	pbUser "jakarta/app/usercenter/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBlockListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBlockListLogic {
	return &GetUserBlockListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserBlockListLogic) GetUserBlockList(req *types.GetUserBlockListReq) (resp *types.GetUserBlockListResp, err error) {
	var in pbUser.GetUserBlockListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.UsercenterRpc.GetUserBlockerList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetUserBlockListResp{List: make([]*types.BlockUserInfo, 0)}
	_ = copier.Copy(resp, rs)
	return
}
