package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"
	"jakarta/app/pgModel/chatPgModel"
	"jakarta/common/key/chatkey"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
)

type UpdateVoiceChatStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateVoiceChatStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVoiceChatStatLogic {
	return &UpdateVoiceChatStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  统计当前通话记录 扣除通话时间
func (l *UpdateVoiceChatStatLogic) UpdateVoiceChatStat(in *pb.UpdateVoiceChatStatReq) (*pb.UpdateVoiceChatStatResp, error) {
	// 判断当前log状态
	data, err := l.svcCtx.VoiceChatLogModel.FindOne(l.ctx, in.ChatLogId)
	if err != nil {
		return nil, err
	}
	if data.State != chatkey.VoiceChatStateStop {
		logx.WithContext(l.ctx).Errorf("UpdateVoiceChatStat state error data:%+v", data)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("wrong log state %d", data.State))
	}
	resp := &pb.UpdateVoiceChatStatResp{}
	state := chatkey.VoiceChatStateSettle
	err = l.svcCtx.ChatBalanceModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var err2 error
		var id string
		// 分别更新用户和XXX当前语音通话状态
		err2 = l.svcCtx.UserVoiceChatStateModel.UpdateState(ctx, session, in.Uid, in.ListenerUid, state)
		if err2 != nil {
			return err2
		}
		err2 = l.svcCtx.ListenerVoiceChatStateModel.UpdateState(ctx, session, in.Uid, in.ListenerUid, state)
		if err2 != nil {
			return err2
		}
		//
		err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, id = l.svcCtx.ChatBalanceModel.UpdateUsedVoiceChatTime(ctx, session, in.Uid, in.ListenerUid, tool.Abs(in.UsedMinute))
		if err2 != nil {
			return err2
		}
		err2 = l.svcCtx.VoiceChatLogModel.UpdateChatLogState(ctx, session, in.ChatLogId, state)
		if err2 != nil {
			return err2
		}
		_, err2 = l.svcCtx.ChatBalanceLogModel.InsertTrans(ctx, session, &chatPgModel.ChatBalanceLog{
			EventType:     chatkey.ChatStatUpdateTypeVoiceChatDecr,
			EventId:       in.ChatLogId,
			Value:         tool.Abs(in.UsedMinute),
			Uid:           in.Uid,
			ListenerUid:   in.ListenerUid,
			ChatBalanceId: id,
			Id:            uniqueid.GenDataId(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
