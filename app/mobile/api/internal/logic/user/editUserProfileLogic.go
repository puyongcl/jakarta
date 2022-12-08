package user

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserProfileLogic {
	return &EditUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserProfileLogic) EditUserProfile(req *types.EditProfileReq) (resp *types.EditProfileResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	if userId != req.Uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.UserUidNotMatch, fmt.Sprintf("%d:%d", userId, req.Uid))
	}
	var pReq usercenter.EditUserProfileReq
	_ = copier.Copy(&pReq, req)
	userInfoResp, err := l.svcCtx.UsercenterRpc.EditUserProfile(l.ctx, &pReq)
	if err != nil {
		return nil, err
	}

	var userInfo types.UserProfile
	_ = copier.Copy(&userInfo, userInfoResp.User)

	return &types.EditProfileResp{UserProfile: &userInfo}, nil
}
