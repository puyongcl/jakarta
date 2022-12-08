package admin

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.AdminLoginReq) (resp *types.AdminLoginResp, err error) {
	if req.AuthType != userkey.UserAuthTypePasswd || req.AuthKey == "" || req.Password == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "param is empty")
	}
	loginResp, err := l.svcCtx.UsercenterRpc.LoginAdmin(l.ctx, &usercenter.LoginReq{
		AuthType: req.AuthType,
		AuthKey:  req.AuthKey,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AdminLoginResp{}
	_ = copier.Copy(resp, loginResp)

	return resp, nil
}
