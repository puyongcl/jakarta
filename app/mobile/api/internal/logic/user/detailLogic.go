package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/app/usercenter/rpc/usercenter"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.GetUserProfileReq) (resp *types.GetUserProfileResp, err error) {
	//
	//userId := ctxdata.GetUidFromCtx(l.ctx)

	userInfoResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserProfileReq{Uid: req.Uid})
	if err != nil {
		return nil, err
	}

	var userInfo types.UserProfile
	_ = copier.Copy(&userInfo, userInfoResp.User)

	return &types.GetUserProfileResp{UserProfile: &userInfo, OpenId: userInfoResp.OpenId}, nil
}
