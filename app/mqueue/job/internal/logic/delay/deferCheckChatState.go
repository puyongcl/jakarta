package delay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
)

// 通话时间使用完后 延时检查是否结算完成 检查文字聊天结束状态

type DeferCheckChatStateHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeferCheckChatStateHandler(svcCtx *svc.ServiceContext) *DeferCheckChatStateHandler {
	return &DeferCheckChatStateHandler{
		svcCtx: svcCtx,
	}
}

// defer check chat state
func (l *DeferCheckChatStateHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p jobtype.DeferCheckChatStatePayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return err
	}

	switch p.OrderType {
	case orderkey.ListenerOrderTypeVoiceChat:
		return l.checkVoiceChat(ctx, &p)
	case orderkey.ListenerOrderTypeTextChat:
		return l.checkTextChat(ctx, &p)
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数错误 %d", p.OrderType))
	}
}

func (l *DeferCheckChatStateHandler) checkVoiceChat(ctx context.Context, p *jobtype.DeferCheckChatStatePayload) error {
	_, err := l.svcCtx.ChatRpc.SyncChatState(ctx, &pbChat.SyncChatStateReq{
		Uid:         p.Uid,
		ListenerUid: p.ListenerUid,
		Action:      chatkey.ChatAction11,
	})
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("DeferCheckChatStateHandler checkVoiceChat done")
	return nil
}

func (l *DeferCheckChatStateHandler) checkTextChat(ctx context.Context, p *jobtype.DeferCheckChatStatePayload) error {
	msg := kqueue.UpdateChatStatMessage{
		Uid:         p.Uid,
		ListenerUid: p.ListenerUid,
		OrderType:   p.OrderType,
	}
	buf, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueUpdateChatStatClient.Push(string(buf))
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("DeferCheckChatStateHandler checkTextChat done")
	return nil
}
