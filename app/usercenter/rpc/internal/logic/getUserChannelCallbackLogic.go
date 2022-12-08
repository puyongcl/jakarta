package logic

import (
	"context"
	"jakarta/app/pgModel/userPgModel"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserChannelCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserChannelCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserChannelCallbackLogic {
	return &GetUserChannelCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取用户渠道回调地址
func (l *GetUserChannelCallbackLogic) GetUserChannelCallback(in *pb.GetUserChannelCallbackReq) (*pb.GetUserChannelCallbackResp, error) {
	rsp, err := l.svcCtx.UserChannelCallbackModel.FindOne(l.ctx, in.Uid)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, err
	}
	resp := &pb.GetUserChannelCallbackResp{}
	if rsp != nil {
		resp.Channel = rsp.Channel
		resp.Cb = rsp.Cb
	}
	return resp, nil
}
