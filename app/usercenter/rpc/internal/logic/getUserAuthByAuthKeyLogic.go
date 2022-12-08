package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAuthByAuthKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByAuthKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByAuthKeyLogic {
	return &GetUserAuthByAuthKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByAuthKeyLogic) GetUserAuthByAuthKey(in *pb.GetUserAuthByAuthKeyReq) (*pb.GetUserAuthByAuthKeyResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthKeyAuthType2(l.ctx, in.AuthKey, in.AuthType)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, err.Error())
	}
	if userAuth == nil {
		return &pb.GetUserAuthByAuthKeyResp{
			UserAuth: &pb.UserAuth{},
		}, nil
	}

	var respUserAuth usercenter.UserAuth
	_ = copier.Copy(&respUserAuth, userAuth)
	respUserAuth.CreateTime = userAuth.CreateTime.Format(db.DateTimeFormat)
	respUserAuth.UpdateTime = userAuth.UpdateTime.Format(db.DateTimeFormat)

	return &pb.GetUserAuthByAuthKeyResp{
		UserAuth: &respUserAuth,
	}, nil
}
