package logic

import (
	"context"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetFreeTextChatCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetFreeTextChatCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetFreeTextChatCntLogic {
	return &ResetFreeTextChatCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  重置免费聊天次数
func (l *ResetFreeTextChatCntLogic) ResetFreeTextChatCnt(in *pb.ResetFreeTextChatCntReq) (*pb.ResetFreeTextChatCntResp, error) {
	err := l.svcCtx.ChatRedis.ResetListenerFreeMsgCnt(l.ctx, in.Uid, in.ListenerUid)
	if err != nil {
		return nil, err
	}
	return &pb.ResetFreeTextChatCntResp{}, nil
}
