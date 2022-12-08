package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserShortProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserShortProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserShortProfileLogic {
	return &GetUserShortProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取用户个人资料
func (l *GetUserShortProfileLogic) GetUserShortProfile(in *pb.GetUserShortProfileReq) (*pb.GetUserShortProfileResp, error) {
	user, err := l.svcCtx.UserProfileModel.FindOne(l.ctx, in.Uid)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("id:%d , err:%v", in.Uid, err))
	}

	if user == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid:%d不存在", in.Uid))
	}

	var respUser usercenter.UserProfile
	_ = copier.Copy(&respUser, user)
	if user.Birthday.Valid {
		respUser.Birthday = user.Birthday.Time.Format(db.DateFormat)
		respUser.Age = tool.GetAge2(user.Birthday.Time)
	}

	return &pb.GetUserShortProfileResp{User: &respUser}, nil
}
