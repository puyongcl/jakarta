package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/listener/mq/internal/svc"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/common/key/listenerkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"strconv"
)

/**
Listening
*/
type UpdateListenerUserStatMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateListenerUserStatMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerUserStatMq {
	return &UpdateListenerUserStatMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateListenerUserStatMq) Consume(_, val string) error {
	var message kqueue.UpdateListenerUserStatMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateListenerUserStatMq->Consume Unmarshal err : %+v , val: %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateListenerUserStatMq->execService err : %+v , val :%s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *UpdateListenerUserStatMq) execService(message *kqueue.UpdateListenerUserStatMessage) (err error) {
	var in pbListener.UpdateListenerUserStatReq
	_ = copier.Copy(&in, message)
	_, err = l.svcCtx.ListenerRpc.UpdateListenerUserStat(l.ctx, &in)
	if err != nil {
		return err
	}
	// 浏览通知
	switch message.Event {
	case listenerkey.ListenerUserEventView:
		err = l.viewNotify(message)
		if err != nil {
			return err
		}
	default:

	}
	return nil
}

// 消息通知
func (l *UpdateListenerUserStatMq) viewNotify(message *kqueue.UpdateListenerUserStatMessage) error {
	for idx := 0; idx < len(message.ListenerUid); idx++ {
		if message.Uid == message.ListenerUid[idx] {
			continue
		}
		kqMsg := kqueue.SendImDefineMessage{
			FromUid:           notify.TimViewNotifyUid,
			ToUid:             message.ListenerUid[idx],
			MsgType:           notify.DefineNotifyMsgTypeViewMsg17,
			Title:             notify.DefineNotifyMsgTemplateViewMsgTitle17,
			Text:              fmt.Sprintf(notify.DefineNotifyMsgTemplateViewMsg17, message.NickName),
			Val1:              strconv.FormatInt(message.Uid, 10),
			Val2:              message.Avatar,
			Val3:              message.NickName,
			Sync:              tim.TimMsgSyncFromNo,
			RepeatMsgCheckId:  fmt.Sprintf("%d-%d", message.Uid, message.ListenerUid[idx]),
			RepeatMsgCheckSec: notify.SendViewNotifyLimitMin * 60,
		}

		// 发送im消息
		var buf []byte
		var err error
		buf, err = json.Marshal(&kqMsg)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			return err
		}
	}
	return nil
}
