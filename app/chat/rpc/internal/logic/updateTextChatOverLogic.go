package logic

import (
	"context"
	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTextChatOverLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTextChatOverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTextChatOverLogic {
	return &UpdateTextChatOverLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新文字聊天时间用完
func (l *UpdateTextChatOverLogic) UpdateTextChatOver(in *pb.UpdateTextChatOverReq) (*pb.UpdateTextChatOverResp, error) {
	err, _, state := l.svcCtx.ChatBalanceModel.UpdateTextChatTimeOver(l.ctx, in.Uid, in.ListenerUid)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTextChatOverResp{State: state}, nil
}
