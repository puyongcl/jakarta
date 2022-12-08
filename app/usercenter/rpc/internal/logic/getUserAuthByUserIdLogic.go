package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAuthByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAuthByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAuthByUserIdLogic {
	return &GetUserAuthByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAuthByUserIdLogic) GetUserAuthByUserId(in *pb.GetUserAuthByUserIdReq) (*pb.GetUserAuthyUserIdResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOne(l.ctx, in.Uid)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("err : %v , in : %+v", err, in))
	}
	if userAuth == nil {
		return &pb.GetUserAuthyUserIdResp{UserAuth: &pb.UserAuth{}}, nil
	}

	var respUserAuth usercenter.UserAuth
	_ = copier.Copy(&respUserAuth, userAuth)
	respUserAuth.CreateTime = userAuth.CreateTime.Format(db.DateTimeFormat)
	respUserAuth.UpdateTime = userAuth.UpdateTime.Format(db.DateTimeFormat)

	return &pb.GetUserAuthyUserIdResp{
		UserAuth: &respUserAuth,
	}, nil
}
