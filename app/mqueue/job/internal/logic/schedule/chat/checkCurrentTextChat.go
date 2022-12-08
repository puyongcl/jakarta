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

type CheckCurrentTextChatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCheckCurrentTextChatHandler(svcCtx *svc.ServiceContext) *CheckCurrentTextChatHandler {
	return &CheckCurrentTextChatHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *CheckCurrentTextChatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	expiryTime := time.Now().Add(time.Duration(chatkey.NotifyOrderOverBeforeMinute) * time.Minute).Format(db.DateTimeFormat)
	var cnt int
	for ; ; pageNo++ {
		var rsp *pbChat.GetUseOutTextChatResp
		var err error
		rsp, err = l.svcCtx.ChatRpc.GetUseOutTextChat(ctx, &pbChat.GetUseOutTextChatReq{
			PageNo:     pageNo,
			PageSize:   10,
			ExpiryTime: expiryTime,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("CheckCurrentTextChatHandler GetUseOutTextChat err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("CheckCurrentTextChatHandler exit. process user cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("CheckCurrentTextChatHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("CheckCurrentTextChatHandler done")
	return nil
}

func (l *CheckCurrentTextChatHandler) update(ctx context.Context, inList []*pbChat.TextChatUser) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		// 发送到时间提醒消息
		msg := &kqueue.SendImDefineMessage{
			FromUid:           inList[idx].ListenerUid,
			ToUid:             inList[idx].Uid,
			MsgType:           notify.DefineNotifyMsgTypeChatMsg21,
			Title:             "",
			Text:              notify.DefineNotifyMsgTemplateChatMsg21,
			Val1:              "1",
			Val2:              "",
			Sync:              tim.TimMsgSyncFromYes,
			RepeatMsgCheckId:  fmt.Sprintf("%d-%d-%d", inList[idx].Uid, inList[idx].ListenerUid, notify.DefineNotifyMsgTypeChatMsg21),
			RepeatMsgCheckSec: chatkey.CheckTextChatUseOutIntervalMinute * chatkey.CheckChatUseOutRange * 2 * 60,
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

		// 延迟更新文字聊天结束状态
		msg2 := &kqueue.CheckChatStateMessage{
			Uid:         inList[idx].Uid,
			ListenerUid: inList[idx].ListenerUid,
			OrderType:   orderkey.ListenerOrderTypeTextChat,
			DeferMinute: chatkey.DeferCheckTextChatMinute,
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
