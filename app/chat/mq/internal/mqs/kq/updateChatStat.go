package kq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/chat/mq/internal/svc"
	pbChat "jakarta/app/chat/rpc/pb"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
	"time"
)

type UpdateChatStatMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateChatStatMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatStatMq {
	return &UpdateChatStatMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateChatStatMq) Consume(_, val string) error {
	var message kqueue.UpdateChatStatMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateChatStatMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}
	logx.WithContext(l.ctx).Infof("UpdateChatStatMq->Consume %s", val)
	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateChatStatMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *UpdateChatStatMq) execService(message *kqueue.UpdateChatStatMessage) error {
	switch message.OrderType {
	case orderkey.ListenerOrderTypeVoiceChat:
		return l.updateVoiceChat(message)
	case orderkey.ListenerOrderTypeTextChat:
		return l.updateTextChat(message)
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "error order type")

	}
}

func (l *UpdateChatStatMq) updateVoiceChat(message *kqueue.UpdateChatStatMessage) error {
	// 计算使用时间
	startTime, err := time.ParseInLocation(db.DateTimeFormat, message.StartTime, time.Local)
	if err != nil {
		return err
	}
	stopTime, err := time.ParseInLocation(db.DateTimeFormat, message.StopTime, time.Local)
	if err != nil {
		return err
	}
	sec := int64(stopTime.Sub(startTime).Seconds())
	vsec := sec % 60
	usedMin := sec / 60
	if vsec > 0 {
		usedMin += 1
	}

	// 更新chat stat
	var rsp *pbChat.UpdateVoiceChatStatResp
	rsp, err = l.svcCtx.ChatRpc.UpdateVoiceChatStat(l.ctx, &pbChat.UpdateVoiceChatStatReq{
		Uid:         message.Uid,
		ListenerUid: message.ListenerUid,
		UsedMinute:  usedMin,
		ChatLogId:   message.LogId,
	})
	if err != nil {
		return err
	}
	var action int64
	if rsp.VoiceChatMinute <= 0 { // 服务时间用完
		action = orderkey.ChatOrderStateUseOutWaitUserConfirm13
	} else {
		action = orderkey.ChatOrderStateUsing4
	}
	// 更新order
	in := pbOrder.UpdateChatOrderUseReq{
		Uid:         message.Uid,
		ListenerUid: message.ListenerUid,
		UsedMinute:  usedMin,
		StartTime:   message.StartTime,
		StopTime:    message.StopTime,
		OrderType:   orderkey.ListenerOrderTypeVoiceChat,
		Action:      action,
	}
	_, err = l.svcCtx.OrderRpc.UpdateChatOrderUse(l.ctx, &in)
	if err != nil {
		return err
	}
	return nil
}

func (l *UpdateChatStatMq) updateTextChat(message *kqueue.UpdateChatStatMessage) error {
	// 更新聊天余额中的文字聊天状态
	in := pbChat.UpdateTextChatOverReq{
		Uid:         message.Uid,
		ListenerUid: message.ListenerUid,
	}
	rsp, err := l.svcCtx.ChatRpc.UpdateTextChatOver(l.ctx, &in)
	if err != nil {
		return err
	}
	if rsp.State != chatkey.TextChatStateStop {
		return nil
	}
	// 更新文字订单状态结束
	in2 := new(pbOrder.UpdateChatOrderUseReq)
	*in2 = pbOrder.UpdateChatOrderUseReq{
		OrderType:   orderkey.ListenerOrderTypeTextChat,
		Uid:         message.Uid,
		ListenerUid: message.ListenerUid,
		Action:      orderkey.ChatOrderStateUseOutWaitUserConfirm13,
	}
	_, err = l.svcCtx.OrderRpc.UpdateChatOrderUse(l.ctx, in2)
	if err != nil {
		return err
	}
	return nil
}
