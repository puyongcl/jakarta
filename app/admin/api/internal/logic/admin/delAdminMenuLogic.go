package admin

import (
	"context"
	"jakarta/app/admin/api/internal/logic/adminlog"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelAdminMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelAdminMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAdminMenuLogic {
	return &DelAdminMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelAdminMenuLogic) DelAdminMenu(req *types.DelAdminMenuReq) (resp *types.DelAdminMenuResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "DelAdminMenu", req.SuperAdminUid, err, req, resp)
	}()
	err = l.svcCtx.AdminMenuPermissionsModel.Delete2(l.ctx, req.Uid, req.Menu1Id, req.Menu2Id)
	if err != nil {
		return
	}
	return &types.DelAdminMenuResp{}, nil
}
