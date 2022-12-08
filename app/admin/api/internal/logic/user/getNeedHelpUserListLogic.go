package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNeedHelpUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNeedHelpUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNeedHelpUserListLogic {
	return &GetNeedHelpUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNeedHelpUserListLogic) GetNeedHelpUserList(req *types.GetNeedHelpUserListReq) (resp *types.GetNeedHelpUserListResp, err error) {
	var in pb.GetNeedHelpUserListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.UsercenterRpc.GetNeedHelpUserList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetNeedHelpUserListResp{
		List: make([]*types.NeedHelpUserData, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
