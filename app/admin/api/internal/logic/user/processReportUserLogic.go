package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/admin/api/internal/logic/adminlog"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/tool"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessReportUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProcessReportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessReportUserLogic {
	return &ProcessReportUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProcessReportUserLogic) ProcessReportUser(req *types.AdminProcessReportUserReq) (resp *types.AdminProcessReportUserResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "ProcessNeedHelpUserLogic", req.AdminUid, err, req, resp)
	}()
	var in pb.AdminProcessReportUserReq
	_ = copier.Copy(&in, req)
	if req.Uid == 0 && tool.IsInt64ArrayExist(req.Action, userkey.AdminBanUserAction) {
		in.Uid = req.AdminUid
	}
	_, err = l.svcCtx.UsercenterRpc.ProcessReportUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.AdminProcessReportUserResp{}
	return
}
