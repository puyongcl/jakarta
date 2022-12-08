package user

import (
	"context"
	"jakarta/common/key/db"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserControlConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserControlConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserControlConfigLogic {
	return &GetUserControlConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserControlConfigLogic) GetUserControlConfig(req *types.GetUserControlConfigReq) (resp *types.GetUserControlConfigResp, err error) {
	resp = &types.GetUserControlConfigResp{}
	resp.LatestAppVer = l.svcCtx.Config.AppVerConf.LatestAppVer
	resp.MinAppVer = l.svcCtx.Config.AppVerConf.MinAppVer
	if req.AppVer <= l.svcCtx.Config.AppVerConf.StoryTabMaxVer {
		resp.StoryTabSwitch = db.Enable
	} else {
		resp.StoryTabSwitch = db.Disable
	}
	return
}
