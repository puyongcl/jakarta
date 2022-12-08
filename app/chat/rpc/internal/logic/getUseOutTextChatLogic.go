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

type GetUseOutTextChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUseOutTextChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUseOutTextChatLogic {
	return &GetUseOutTextChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取时间快结束的文字通话
func (l *GetUseOutTextChatLogic) GetUseOutTextChat(in *pb.GetUseOutTextChatReq) (*pb.GetUseOutTextChatResp, error) {
	tt, err := time.ParseInLocation(db.DateTimeFormat, in.ExpiryTime, time.Local)
	if err != nil {
		return nil, err
	}
	expiryStart := tt.Add(-time.Duration(chatkey.CheckTextChatUseOutIntervalMinute) * time.Minute)
	expiryEnd := tt
	rs, err := l.svcCtx.ChatBalanceModelRO.FindTextChat(l.ctx, 0, 0, &expiryStart, &expiryEnd, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetUseOutTextChatResp{List: make([]*pb.TextChatUser, 0)}
	if len(rs) <= 0 {
		return resp, nil
	}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.TextChatUser
		val.Uid = rs[idx].Uid
		val.ListenerUid = rs[idx].ListenerUid
		val.CurrentChatLogId = rs[idx].CurrentChatLogId
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
