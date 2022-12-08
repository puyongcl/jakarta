package logic

import (
	"context"
	chatPgModel2 "jakarta/app/pgModel/chatPgModel"
	"jakarta/common/key/chatkey"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateListenerVoiceChatStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateListenerVoiceChatStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateListenerVoiceChatStateLogic {
	return &CreateListenerVoiceChatStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  初始化XXX通话状态
func (l *CreateListenerVoiceChatStateLogic) CreateListenerVoiceChatState(in *pb.CreateListenerVoiceChatStateReq) (*pb.CreateListenerVoiceChatStateResp, error) {
	vcsl, err := l.svcCtx.ListenerVoiceChatStateModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != chatPgModel2.ErrNotFound {
		return nil, err
	}
	if vcsl == nil {
		_, err = l.svcCtx.ListenerVoiceChatStateModel.Insert(l.ctx, &chatPgModel2.ListenerVoiceChatState{
			ListenerUid: in.ListenerUid,
			State:       chatkey.VoiceChatStateSettle,
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.CreateListenerVoiceChatStateResp{}, nil
}
