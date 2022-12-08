package logic

import (
	"context"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginAdminLogic {
	return &LoginAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginAdminLogic) LoginAdmin(in *pb.LoginReq) (*pb.LoginResp, error) {
	var err error
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthKeyAuthType2(l.ctx, in.AuthKey, in.AuthType)
	if err != nil && err != userPgModel.ErrNotFound {
		return &pb.LoginResp{}, xerr.NewGrpcErrCodeMsg(xerr.DbError, err.Error())
	}
	if userAuth == nil { // not register
		return nil, xerr.NewGrpcErrCodeMsg(xerr.UserNotExist, "admin not exist")
	}
	if userAuth.UserType != userkey.UserTypeAdmin {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "not a admin user")
	}

	return doLogin(l.ctx, l.svcCtx, in, userAuth)
}
