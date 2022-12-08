package logic

import (
	"context"
	"fmt"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLoginStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLoginStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLoginStateLogic {
	return &UpdateUserLoginStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLoginStateLogic) UpdateUserLoginState(in *pb.UpdateUserLoginStateReq) (*pb.UpdateUserLoginStateResp, error) {
	// 用户登陆状态
	uls, err := l.svcCtx.UserLoginStateModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("not find login state uid:%d err:%+v", in.Uid, err))
	}

	eventTime := time.Now()
	if in.EventTime > 0 {
		// 1665976726538
		eventTime = time.UnixMilli(in.EventTime)
	}

	switch in.State {
	case userkey.Login:
		//
		if uls.LoginTime.Day() != eventTime.Day() { // 今日首次登陆
			uls.LoginCntToday = 1

		} else {
			uls.LoginCntToday++
		}

		uls.LoginCntSum++
		uls.LoginTime = eventTime

	case userkey.Logout, userkey.Disconnect:
		uls.OfflineTime = eventTime
	default:
		logx.WithContext(l.ctx).Errorf("UpdateUserLoginStateLogic uid:%d state:%d", in.Uid, in.State)
		uls.OfflineTime = eventTime
	}

	uls.LoginState = in.State
	uls.ImEventTime = eventTime

	var ra int64
	ra, err = l.svcCtx.UserLoginStateModel.UpdateLoginState(l.ctx, uls)
	if err != nil {
		return nil, err
	}

	//
	auth, err := l.svcCtx.UserAuthModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	logx.WithContext(l.ctx).Infof("UpdateUserLoginStateLogic UID:%d LOGIN STATE:%d USER TYPE:%d Ra:%d val:%+v", in.Uid, in.State, auth.UserType, ra, uls)
	return &pb.UpdateUserLoginStateResp{UserType: auth.UserType, Channel: auth.Channel, TodayLoginCnt: uls.LoginCntToday, IsUpdated: ra}, nil
}
