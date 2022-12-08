package logic

import (
	"context"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"time"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUseOutVoiceChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUseOutVoiceChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUseOutVoiceChatLogic {
	return &GetUseOutVoiceChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取时间快结束的语音通话
func (l *GetUseOutVoiceChatLogic) GetUseOutVoiceChat(in *pb.GetUseOutVoiceChatReq) (*pb.GetUseOutVoiceChatResp, error) {
	tt, err := time.ParseInLocation(db.DateTimeFormat, in.ExpiryTime, time.Local)
	if err != nil {
		return nil, err
	}
	expiryStart := tt.Add(-time.Duration(chatkey.CheckChatUseOutRange*chatkey.CheckVoiceChatUseOutIntervalSecond) * time.Second)
	expiryEnd := tt.Add(time.Duration(chatkey.CheckVoiceChatUseOutIntervalSecond) * time.Second)
	rs, err := l.svcCtx.ChatBalanceModelRO.FindVoiceChat(l.ctx, 0, 0, in.State, &expiryStart, &expiryEnd, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetUseOutVoiceChatResp{List: make([]*pb.VoiceChatUser, 0)}
	if len(rs) <= 0 {
		return resp, nil
	}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.VoiceChatUser
		val.Uid = rs[idx].Uid
		val.ListenerUid = rs[idx].ListenerUid
		val.CurrentChatLogId = rs[idx].CurrentChatLogId
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
