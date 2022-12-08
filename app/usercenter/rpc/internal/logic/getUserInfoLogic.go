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
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNoExistsError = xerr.NewErrCode(xerr.UserNotExist)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserProfileReq) (*pb.GetUserProfileResp, error) {
	//
	user, err := l.svcCtx.UserProfileModel.FindOne(l.ctx, in.Uid)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("id:%d , err:%v", in.Uid, err))
	}

	if user == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid:%d不存在", in.Uid))
	}
	//
	userauth, err := l.svcCtx.UserAuthModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	var respUser usercenter.UserProfile
	_ = copier.Copy(&respUser, user)
	if user.Birthday.Valid {
		respUser.Birthday = user.Birthday.Time.Format(db.DateFormat)
		respUser.Age = tool.GetAge2(user.Birthday.Time)
	}
	respUser.CreateTime = user.CreateTime.Format(db.DateTimeFormat)

	return &pb.GetUserProfileResp{User: &respUser, OpenId: userauth.AuthKey}, nil
}
