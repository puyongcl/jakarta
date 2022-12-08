package logic

import (
	"context"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOnlineStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserOnlineStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOnlineStateLogic {
	return &GetUserOnlineStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取用户在线状态和用户类型
func (l *GetUserOnlineStateLogic) GetUserOnlineState(in *pb.GetUserOnlineStateReq) (*pb.GetUserOnlineStateResp, error) {
	state, err := l.svcCtx.UserLoginStateModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserOnlineStateResp{Uid: in.Uid, OnlineState: state.LoginState}, nil
}
