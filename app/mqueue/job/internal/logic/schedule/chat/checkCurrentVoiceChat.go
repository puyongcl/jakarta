package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"time"
)

type CheckCurrentVoiceChatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCheckCurrentVoiceChatHandler(svcCtx *svc.ServiceContext) *CheckCurrentVoiceChatHandler {
	return &CheckCurrentVoiceChatHandler{
		svcCtx: svcCtx,
	}
}

// 1分钟结束
func (l *CheckCurrentVoiceChatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	expiryTime := time.Now().Add(time.Duration(chatkey.NotifyOrderOverBeforeMinute) * time.Minute).Format(db.DateTimeFormat)
	var cnt int
	for ; ; pageNo++ {
		var rsp *pbChat.GetUseOutVoiceChatResp
		var err error
		rsp, err = l.svcCtx.ChatRpc.GetUseOutVoiceChat(ctx, &pbChat.GetUseOutVoiceChatReq{
			PageNo:     pageNo,
			PageSize:   10,
			ExpiryTime: expiryTime,
			State:      chatkey.VoiceChatStateStart,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("CheckCurrentVoiceChatHandler GetUseOutVoiceChat err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("CheckCurrentVoiceChatHandler exit. process user cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("CheckCurrentVoiceChatHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("CheckCurrentVoiceChatHandler done")
	return nil
}

func (l *CheckCurrentVoiceChatHandler) update(ctx context.Context, inList []*pbChat.VoiceChatUser) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		// 发送到时间提醒消息
		msg := &kqueue.SendImDefineMessage{
			FromUid:           inList[idx].ListenerUid,
			ToUid:             inList[idx].Uid,
			MsgType:           notify.DefineNotifyMsgTypeChatMsg20,
			Title:             "",
			Text:              notify.DefineNotifyMsgTemplateChatMsg20,
			Val1:              "1",
			Val2:              "",
			Sync:              tim.TimMsgSyncFromYes,
			RepeatMsgCheckId:  fmt.Sprintf("%d-%d-%s", inList[idx].Uid, inList[idx].ListenerUid, inList[idx].CurrentChatLogId),
			RepeatMsgCheckSec: chatkey.CheckVoiceChatUseOutIntervalSecond * 18,
		}
		var buf []byte
		buf, err = json.Marshal(msg)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			return err
		}

		// 延迟检查是否结算
		msg2 := &kqueue.CheckChatStateMessage{
			Uid:         inList[idx].Uid,
			ListenerUid: inList[idx].ListenerUid,
			OrderType:   orderkey.ListenerOrderTypeVoiceChat,
			DeferMinute: chatkey.DeferCheckVoiceChatMinute,
		}
		var buf2 []byte
		buf2, err = json.Marshal(msg2)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueCheckChatStateClient.Push(string(buf2))
		if err != nil {
			return err
		}
	}
	return
}
