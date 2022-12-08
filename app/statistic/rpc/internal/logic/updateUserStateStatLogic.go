package logic

import (
	"context"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/userkey"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStateStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStateStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStateStatLogic {
	return &UpdateUserStateStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  定时统计用户和XXX状态数据
func (l *UpdateUserStateStatLogic) UpdateUserStateStat(in *pb.UpdateUserStateStatReq) (*pb.UpdateUserStateStatResp, error) {
	// 统计接单状态
	var s2c, s3c, s4c int64
	var err error
	resp := &pb.UpdateUserStateStatResp{}
	s2c, err = l.svcCtx.ListenerProfileModel.CountListenerByWorkState(l.ctx, listenerkey.ListenerWorkStateWorking)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountListenerByWorkState err:%+v", err)
		return resp, nil
	}
	s3c, err = l.svcCtx.ListenerProfileModel.CountListenerByWorkState(l.ctx, listenerkey.ListenerWorkStateRestingAuto)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountListenerByWorkState err:%+v", err)
		return resp, nil
	}
	s4c, err = l.svcCtx.ListenerProfileModel.CountListenerByWorkState(l.ctx, listenerkey.ListenerWorkStateRestingManual)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountListenerByWorkState err:%+v", err)
		return resp, nil
	}
	// 统计在线状态
	var ol1, ol2, ou1 int64
	ol1, err = l.svcCtx.ListenerProfileModel.CountListenerByOnlineState(l.ctx, userkey.Login)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountListenerByOnlineState err:%+v", err)
		return resp, nil
	}

	ol2, err = l.svcCtx.ListenerProfileModel.CountListenerByOnlineState(l.ctx, 0)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountListenerByOnlineState err:%+v", err)
		return resp, nil
	}

	// 统计用户在线状态
	ou1, err = l.svcCtx.UserLoginStateModel.CountUserOnline(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic CountUserOnline err:%+v", err)
		return resp, nil
	}

	logx.WithContext(l.ctx).Infof("UpdateUserStateStatLogic Listener Work State Stat Working:%d RestingAuto:%d RestingManual:%d Online:%d Offline:%d User Online:%d", s2c, s3c, s4c, ol1, ol2, ou1)

	return &pb.UpdateUserStateStatResp{}, nil
}
