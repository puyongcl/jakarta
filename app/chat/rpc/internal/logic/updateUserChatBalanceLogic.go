package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	chatPgModel2 "jakarta/app/pgModel/chatPgModel"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserChatBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserChatBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserChatBalanceLogic {
	return &UpdateUserChatBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新用户聊天可用时间
func (l *UpdateUserChatBalanceLogic) UpdateUserChatBalance(in *pb.UpdateUserChatBalanceReq) (*pb.UpdateUserChatBalanceResp, error) {
	if in.AddMinute == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	resp := &pb.UpdateUserChatBalanceResp{}
	err := l.svcCtx.ChatBalanceModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var err2 error
		var id string
		switch in.EventType {
		case chatkey.ChatStatUpdateTypeOrderPaidAdd:
			if in.AddMinute <= 0 {
				return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
			}
			if in.OrderType == orderkey.ListenerOrderTypeTextChat { // 增加文字聊天时间
				var exTime time.Time
				exTime, err2 = time.ParseInLocation(db.DateTimeFormat, in.TextExpireTime, time.Local)
				if err2 != nil {
					return err2
				}
				err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, id = l.svcCtx.ChatBalanceModel.AddTextChatTime(ctx, session, in.Uid, in.ListenerUid, in.AddMinute, &exTime)
				if err2 != nil {
					return err2
				}
			} else if in.OrderType == orderkey.ListenerOrderTypeVoiceChat { // 增加语音聊天时间
				err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, id = l.svcCtx.ChatBalanceModel.AddVoiceChatTime(ctx, session, in.Uid, in.ListenerUid, in.AddMinute)
				if err2 != nil {
					return err2
				}
			}
		case chatkey.ChatStatUpdateTypeOrderExpireDecr, chatkey.ChatStatUpdateTypeOrderUserStop:
			if in.OrderType == orderkey.ListenerOrderTypeVoiceChat { // 减少语音聊天时间
				err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, id = l.svcCtx.ChatBalanceModel.ReduceVoiceChatTime(ctx, session, in.Uid, in.ListenerUid, tool.Abs(in.AddMinute))
				if err2 != nil {
					return err2
				}
			} else if in.OrderType == orderkey.ListenerOrderTypeTextChat { // 结束文字聊天
				err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, id = l.svcCtx.ChatBalanceModel.StopTextChatTime(ctx, session, in.Uid, in.ListenerUid, tool.Abs(in.AddMinute))
				if err2 != nil {
					return err2
				}
			}

		default:
			return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
		}

		// 余额变化日志
		_, err2 = l.svcCtx.ChatBalanceLogModel.InsertTrans(ctx, session, &chatPgModel2.ChatBalanceLog{
			EventType:     in.EventType,
			EventId:       in.OrderId,
			Value:         in.AddMinute,
			Uid:           in.Uid,
			ListenerUid:   in.ListenerUid,
			ChatBalanceId: id,
			Id:            uniqueid.GenDataId(),
		})
		if err2 != nil {
			return err2
		}

		// 检查是否存在用户当前语音通话数据
		vcsu, err2 := l.svcCtx.UserVoiceChatStateModel.FindOne(ctx, in.Uid)
		if err2 != nil && err2 != chatPgModel2.ErrNotFound {
			return err2
		}
		if vcsu == nil {
			_, err2 = l.svcCtx.UserVoiceChatStateModel.InsertTrans(ctx, session, &chatPgModel2.UserVoiceChatState{
				Uid: in.Uid,
			})
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
