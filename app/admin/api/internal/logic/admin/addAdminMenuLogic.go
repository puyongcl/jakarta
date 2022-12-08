package admin

import (
	"context"
	"encoding/json"
	"jakarta/app/admin/api/internal/logic/adminlog"
	"jakarta/app/pgModel/adminPgModel"
	"jakarta/common/key/db"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAdminMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAdminMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminMenuLogic {
	return &AddAdminMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAdminMenuLogic) AddAdminMenu(req *types.AddAdminMenuReq) (resp *types.AddAdminMenuResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "AddAdminMenu", req.SuperAdminUid, err, req, resp)
	}()
	data := new(adminPgModel.AdminMenuPermissions)
	data.Menu1Id = req.Menu1Id
	data.Menu2Id = req.Menu2Id
	data.State = db.Enable
	data.Uid = req.Uid
	js, err := json.Marshal(req.MenuValue)
	if err != nil {
		return nil, err
	}
	data.MenuValue = string(js)

	_, err = l.svcCtx.AdminMenuPermissionsModel.Insert2(l.ctx, data)
	if err != nil {
		return nil, err
	}
	return
}
