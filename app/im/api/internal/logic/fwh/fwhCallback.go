package fwh

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/api/internal/svc"
	"jakarta/common/key/fwh"
	"jakarta/common/kqueue"
	"net/http"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) WxFwhCallback(w http.ResponseWriter, r *http.Request) (err error) {
	wxsrv := l.svcCtx.Wxfwh.GetServer(r, w)
	var event message.EventType
	wxsrv.SetMessageHandler(func(mixMsg *message.MixMessage) *message.Reply {
		event = mixMsg.Event
		switch mixMsg.Event {
		case message.EventSubscribe:
			txt := message.NewText(fmt.Sprintf(fwh.SubscribeReplyMsg, l.svcCtx.Config.WxMiniConf.AppId))
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: txt}
		default:
			return nil
		}
	})
	err = wxsrv.Serve()
	if err != nil {
		return
	}

	switch event {
	case message.EventSubscribe, message.EventUnsubscribe: // 关注和取消关注事件
		err = l.userSubscribeEvent(wxsrv.GetOpenID(), event)
		if err != nil {
			return err
		}

	default:

	}

	if wxsrv.ResponseMsg != nil { // 发送回复消息
		return wxsrv.Send()
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("success"))
	return
}

// 关注事件
func (l *CallbackLogic) userSubscribeEvent(openId string, event message.EventType) error {
	kqMsg := kqueue.WxFwhCallbackEventMessage{
		OpenId: openId,
		Event:  string(event),
	}
	buf, err := json.Marshal(&kqMsg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueWxFwhCallbackEventClient.Push(string(buf))
	if err != nil {
		return err
	}
	return nil
}
