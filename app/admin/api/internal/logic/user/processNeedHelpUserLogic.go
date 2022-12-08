package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/admin/api/internal/logic/adminlog"
	"jakarta/app/usercenter/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessNeedHelpUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessNeedHelpUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessNeedHelpUserLogic {
	return &ProcessNeedHelpUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessNeedHelpUserLogic) ProcessNeedHelpUser(req *types.ProcessNeedHelpUserReq) (resp *types.ProcessNeedHelpUserResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "ProcessNeedHelpUserLogic", req.AdminUid, err, req, resp)
	}()
	var in pb.AdminMarkNeedHelpUserReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.UsercenterRpc.ProcessNeedHelpUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.ProcessNeedHelpUserResp{}
	return
}
