package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReportUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReportUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReportUserListLogic {
	return &GetReportUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReportUserListLogic) GetReportUserList(req *types.GetReportUserListReq) (resp *types.GetReportUserListResp, err error) {
	var in pb.GetReportUserListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.UsercenterRpc.GetReportUserList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetReportUserListResp{
		List: make([]*types.ReportUserData, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
