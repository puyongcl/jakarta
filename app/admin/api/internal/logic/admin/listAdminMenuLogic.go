package admin

import (
	"context"
	"encoding/json"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdminMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAdminMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdminMenuLogic {
	return &ListAdminMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAdminMenuLogic) ListAdminMenu(req *types.ListAdminMenuReq) (resp *types.ListAdminMenuResp, err error) {
	resp = &types.ListAdminMenuResp{
		List: make([]*types.AdminMenu, 0),
		Sum:  0,
	}
	resp.Sum, err = l.svcCtx.AdminMenuPermissionsModel.FindCount(l.ctx, req.Uid, req.Menu1Id, req.Menu2Id)
	if err != nil {
		return nil, err
	}
	if resp.Sum <= 0 {
		return
	}
	rsp, err := l.svcCtx.AdminMenuPermissionsModel.Find(l.ctx, req.Uid, req.Menu1Id, req.Menu2Id, req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	for idx := 0; idx < len(rsp); idx++ {
		var menu types.AdminMenu
		err = json.Unmarshal([]byte(rsp[idx].MenuValue), &menu)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, &menu)
	}
	return
}
