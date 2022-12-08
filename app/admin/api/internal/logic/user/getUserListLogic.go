package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/app/statistic/rpc/pb"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListResp, err error) {
	var in pb.GetUserListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.StatRpc.GetUserList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetUserListResp{
		List: make([]*types.UserDetail, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
