package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/chat/mq/internal/svc"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/jobtype"
	pbOrder "jakarta/app/order/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"strconv"
	"strings"
	"time"
)

type UserEnterChatMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserEnterChatMq(ctx context.Context, svcCtx *svc.ServiceContext) *UserEnterChatMq {
	return &UserEnterChatMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserEnterChatMq) Consume(_, val string) error {
	var message kqueue.UserFirstEnterChatMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *UserEnterChatMq) execService(msg *kqueue.UserFirstEnterChatMessage) error {
	if msg.IsFirst == db.Enable { // 首次进入聊天
		// 发送个人介绍
		in := pbListener.GetListenerBasicInfoReq{ListenerUid: msg.ListenerUid}
		rsp, err := l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &in)
		if err != nil {
			return err
		}
		//
		l.sendNotify(msg.ListenerUid, msg.Uid, rsp.Introduction)
	}
	// 发送氛围提示语 两天内不再发送
	l.sendNotify2(msg.ListenerUid, msg.Uid, msg.IsFirst)

	// 统计浏览用户 发送浏览通知
	l.sendViewMsg(msg)

	return nil
}

func (l *UserEnterChatMq) sendViewMsg(msg *kqueue.UserFirstEnterChatMessage) {
	var in pbListener.UpdateListenerUserStatReq
	in.Uid = msg.Uid
	in.Event = listenerkey.ListenerUserEventView
	in.Time = time.Now().Format(db.DateTimeFormat)
	in.ListenerUid = []int64{msg.ListenerUid}
	_, err := l.svcCtx.ListenerRpc.UpdateListenerUserStat(l.ctx, &in)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendViewMsg UpdateListenerUserStat err:%+v", err)
		return
	}

	in2 := pbUser.GetUserShortProfileReq{Uid: msg.Uid}
	rsp2, err := l.svcCtx.UsercenterRpc.GetUserShortProfile(l.ctx, &in2)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendViewMsg GetUserShortProfile err:%+v", err)
		return
	}
	kqMsg := kqueue.SendImDefineMessage{
		FromUid:           notify.TimViewNotifyUid,
		ToUid:             msg.ListenerUid,
		MsgType:           notify.DefineNotifyMsgTypeViewMsg17,
		Title:             notify.DefineNotifyMsgTemplateViewMsgTitle17,
		Text:              fmt.Sprintf(notify.DefineNotifyMsgTemplateViewMsg17, rsp2.User.Nickname),
		Val1:              strconv.FormatInt(msg.Uid, 10),
		Val2:              rsp2.User.Avatar,
		Val3:              rsp2.User.Nickname,
		Sync:              tim.TimMsgSyncFromNo,
		RepeatMsgCheckId:  fmt.Sprintf("%d-%d", msg.Uid, msg.ListenerUid),
		RepeatMsgCheckSec: notify.SendViewNotifyLimitMin * 60,
	}

	// 发送im消息
	var buf []byte
	buf, err = json.Marshal(&kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendViewMsg marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendViewMsg Push err:%+v", err)
		return
	}
}

func (l *UserEnterChatMq) sendNotify2(listenerUid, uid, isFirst int64) {
	id := fmt.Sprintf("%d-%d-%d", listenerUid, uid, notify.DefineNotifyMsgTypeChatMsg31)
	b, err := l.svcCtx.IMRedis.IsHaveNotifyLimit(l.ctx, id, notify.DefineNotifyMsgTypeChatMsg31)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify2 IsHaveNotifyLimit err:%+v", err)
		return
	}

	if b {
		return
	}
	var lp *pbListener.GetListenerBasicInfoResp
	var in pbListener.GetListenerBasicInfoReq
	in.ListenerUid = listenerUid
	lp, err = l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &in)
	if err != nil {
		return
	}

	var txt string
	now := time.Now()
	if (int64(now.Minute())+listenerUid+uid)%2 == 0 { // 发送收到满意评价
		var in2 pbOrder.GetLastCommentOrderReq
		in2.Star = listenerkey.Rating5Star
		in2.ListenerUid = listenerUid
		var rsp *pbOrder.GetLastCommentOrderResp
		rsp, err = l.svcCtx.OrderRpc.GetLastCommentOrder(l.ctx, &in2)
		if err != nil {
			return
		}
		var ct time.Time
		ct, err = time.ParseInLocation(db.DateTimeFormat, rsp.CommentTime, time.Local)
		if err != nil {
			return
		}
		ts := tool.GetTimeDurationText(now.Sub(ct))
		txt = fmt.Sprintf(notify.DefineNotifyMsgTemplateChatMsg31A, strings.Replace(lp.NickName, "#", "", -1), ts)
	} else { // 帮助多少人
		var in2 pbOrder.GetRecentPaidUserCntReq
		in2.ListenerUid = listenerUid
		in2.StartTime = now.AddDate(0, 0, -notify.RecentHelpUserDay).Format(db.DateTimeFormat)
		in2.EndTime = now.Format(db.DateTimeFormat)
		var rsp *pbOrder.GetRecentPaidUserCntResp
		rsp, err = l.svcCtx.OrderRpc.GetRecentPaidUserCnt(l.ctx, &in2)
		if err != nil {
			return
		}
		randCnt := tool.DivideInt64(24, 3) * notify.RecentHelpUserDay
		randCnt = tool.DivideInt64(randCnt, 2)
		addCnt := (uid+listenerUid)%randCnt + randCnt
		txt = fmt.Sprintf(notify.DefineNotifyMsgTemplateChatMsg31B, strings.Replace(lp.NickName, "#", "", -1), notify.RecentHelpUserDay, rsp.UserCnt+addCnt)
	}
	kqm := kqueue.SendImDefineMessage{
		FromUid:           listenerUid,
		ToUid:             uid,
		MsgType:           notify.DefineNotifyMsgTypeChatMsg31,
		Title:             "",
		Text:              txt,
		Val1:              "",
		Val2:              "",
		Val3:              "",
		Val4:              "",
		Val5:              "",
		Val6:              "",
		Sync:              tim.TimMsgSyncFromNo,
		RepeatMsgCheckId:  fmt.Sprintf("%d-%d-%d", listenerUid, uid, notify.DefineNotifyMsgTypeChatMsg31),
		RepeatMsgCheckSec: notify.SendChatMsg31IntervalSecond,
	}

	buf, err := json.Marshal(&kqm)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify2 json Marshal err:%+v", err)
		return
	}

	if isFirst != db.Enable { // 非首次 不发个人介绍的情况下 不做延时发送处理
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify2 kafka Push err:%+v", err)
			return
		}
		return
	}

	// 第一次和个人介绍一起发送时 延迟发送保证顺序
	var payload []byte
	payload, err = json.Marshal(jobtype.DeferSendImMsgPayload{KqMsgBuf: buf})
	if err != nil {
		return
	}

	_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferSendImMsg, payload), asynq.ProcessIn(time.Duration(notify.DeferSendChatMsg31Second)*time.Second))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify EnqueueContext err:%+v", err)
		return
	}
	return
}

// 发送首次进入聊天界面消息
func (l *UserEnterChatMq) sendNotify(listenerUid, uid int64, intro string) {
	kqm := kqueue.SendImDefineMessage{
		FromUid:           listenerUid,
		ToUid:             uid,
		MsgType:           notify.TextMsgTypeIntroMsg1,
		Title:             "",
		Text:              intro,
		Val1:              "",
		Val2:              "",
		Val3:              "",
		Val4:              "",
		Val5:              "",
		Val6:              "",
		Sync:              tim.TimMsgSyncFromNo,
		RepeatMsgCheckId:  fmt.Sprintf("%d-%d-%d", listenerUid, uid, notify.TextMsgTypeIntroMsg1),
		RepeatMsgCheckSec: notify.SendListenerIntroMsgIntervalSecond,
	}
	buf, err := json.Marshal(&kqm)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify json Marshal err:%+v", err)
		return
	}

	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UserEnterChatMq sendNotify kafka Push err:%+v", err)
		return
	}
	return
}
