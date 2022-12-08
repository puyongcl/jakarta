package kq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/xerr"
	"strconv"
)

// 订阅消息
type SubscribeNotifyMsgMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeNotifyMsgMq(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeNotifyMsgMq {
	return &SubscribeNotifyMsgMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeNotifyMsgMq) Consume(_, val string) error {
	var message kqueue.SubscribeNotifyMsgMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("SubscribeNotifyMsgMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("SubscribeNotifyMsgMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *SubscribeNotifyMsgMq) execService(kqMsg *kqueue.SubscribeNotifyMsgMessage) error {
	switch kqMsg.Action {
	case notify.SubscribeUserNotifyMsgEventAdd: // 订阅某人的一种消息
		return l.svcCtx.IMRedis.AddSubscribeNotifyMember(l.ctx, kqMsg.TargetUid, kqMsg.MsgType, kqMsg.SendCnt, kqMsg.Uid)
	case notify.SubscribeUserNotifyMsgEventCancel: // 取消订阅某人的一种消息
		_, err := l.svcCtx.IMRedis.RemSubscribeNotifyMember(l.ctx, kqMsg.TargetUid, kqMsg.MsgType, kqMsg.SendCnt, strconv.FormatInt(kqMsg.Uid, 10))
		if err != nil {
			return err
		}
		return nil
	case notify.SubscribeUserNotifyMsgEventSend: // 发送某人的一种订阅消息
		return l.sendNotifyMsg(kqMsg)
	case notify.SubscribeOneTimeNotifyMsgEventAdd: // 一次性订阅某种消息
		return l.subscribeOneTimeMsg(kqMsg)
	case notify.SubscribeOneTimeNotifyMsgEventSend: // 发送一次性订阅某种消息
		return l.sendSubscribeOneTimeMsg(kqMsg)
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "错误参数")
	}
}

// 一次性订阅消息
func (l *SubscribeNotifyMsgMq) subscribeOneTimeMsg(kqMsg *kqueue.SubscribeNotifyMsgMessage) error {
	return l.svcCtx.IMRedis.AddOneTimeSubscribeMsgMember(l.ctx, kqMsg.MsgType, kqMsg.Uid)
}

// 发送一次性订阅消息
func (l *SubscribeNotifyMsgMq) sendSubscribeOneTimeMsg(kqMsg *kqueue.SubscribeNotifyMsgMessage) error {
	// im消息
	if kqMsg.IMMsg != nil {
		buf, err := json.Marshal(kqMsg.IMMsg)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			return err
		}
	}
	// 小程序订阅通知
	if kqMsg.MpMsg != nil {
		// 检查是否有通知次数
		r, err := l.svcCtx.IMRedis.CostOneTimeSubscribeMsg(l.ctx, kqMsg.MsgType, kqMsg.Uid)
		if err != nil {
			return err
		}
		if !r { // 返回false表示没有订阅
			return nil
		}
		var buf []byte
		buf, err = json.Marshal(kqMsg.MpMsg)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueSendWxMiniProgramMsgClient.Push(string(buf))
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *SubscribeNotifyMsgMq) sendNotifyMsg(message *kqueue.SubscribeNotifyMsgMessage) error {
	for {
		rs, err := l.svcCtx.IMRedis.GetSubscribeNotifyMember(l.ctx, message.TargetUid, message.MsgType, message.SendCnt, 20)
		if err != nil {
			return err
		}
		if len(rs) <= 0 {
			return nil
		}

		for k, _ := range rs {
			message.IMMsg.ToUid, err = strconv.ParseInt(rs[k], 10, 64)
			if err != nil {
				return err
			}
			var buf []byte
			buf, err = json.Marshal(message.IMMsg)
			if err != nil {
				return err
			}
			err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
			if err != nil {
				return err
			}
		}

		// 如果是一次订阅消息 需要移除订阅者
		if message.SendCnt == notify.SubscribeNotifyMsgSendCntOne {
			_, err = l.svcCtx.IMRedis.RemSubscribeNotifyMember(l.ctx, message.TargetUid, message.MsgType, message.SendCnt, rs...)
			if err != nil {
				return err
			}
		}
	}
}
