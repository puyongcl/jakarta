package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/listener/mq/internal/svc"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/jobtype"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"strings"
	"time"
)

// 用户登陆后推荐XXX发消息
type SendHelloWhenUserLoginMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendHelloWhenUserLoginMq(ctx context.Context, svcCtx *svc.ServiceContext) *SendHelloWhenUserLoginMq {
	return &SendHelloWhenUserLoginMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendHelloWhenUserLoginMq) Consume(_, val string) error {
	var message kqueue.SendHelloWhenUserLoginMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendHelloWhenUserLoginMq->Consume Unmarshal err : %+v , val: %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendHelloWhenUserLoginMq->execService err : %+v , val :%s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *SendHelloWhenUserLoginMq) execService(msg *kqueue.SendHelloWhenUserLoginMessage) (err error) {
	return l.doSendHello(msg.Uid, msg.IsNewUser)
}

func (l *SendHelloWhenUserLoginMq) doSendHello(uid int64, isNewUser int64) (err error) {
	// 推荐几位XXX
	var cnt int
	tz := time.Now().Minute()
	cnt = tz%2 + 2

	var listenerUids []int64
	listenerUids, err = l.svcCtx.ListenerRedis.GetRecommendListenerAndIncrScore(l.ctx, rediskey.RedisKeyListenerRecommendSendHelloMsgWhenUserLogin, cnt)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic ListenerRedis.GetRecommendListenerAndIncrScore err:%+v", err)
		return
	}

	// 延迟发送消息
	if len(listenerUids) > 0 {
		if isNewUser == db.Enable {
			l.sendNewUserMsg(uid, listenerUids)
		} else {
			l.sendUserMsg(uid, listenerUids)
		}
	}
	return
}

// 取几条消息 剩下到存到set
func (l *SendHelloWhenUserLoginMq) getSetHelloMsg(uid int64, cnt int64) ([]int64, error) {
	hml := int64(len(notify.LoginHelloMsg2))
	idx := uid % (hml - cnt)

	var rv []int64
	var idx2 int64
	for idx2 = 0; idx2 < cnt; idx2++ {
		rv = append(rv, idx+idx2)
	}

	var rv2 []interface{}
	for idx2 = 0; idx2 < hml; idx2++ {
		if tool.IsInt64ArrayExist(idx2, rv) {
			continue
		}

		rv2 = append(rv2, idx2)
	}
	//
	err := l.svcCtx.ListenerRedis.SetUserHelloMsg(l.ctx, uid, rv2)
	if err != nil {
		return []int64{}, err
	}

	return rv, err
}

// 第一次登陆发送的消息
func (l *SendHelloWhenUserLoginMq) sendNewUserMsg(uid int64, listenerUids []int64) {
	if len(listenerUids) <= 0 {
		return
	}

	// 随机取两条消息 把剩下的存到set
	var err error
	var hm []int64 // 本次要发送到问候消息
	hm, err = l.getSetHelloMsg(uid, int64(len(listenerUids)))
	if err != nil {
		return
	}

	msgs := make([]*kqueue.SendImDefineMessage, len(listenerUids))
	for idx := 0; idx < len(listenerUids); idx++ {
		kqm2 := kqueue.SendImDefineMessage{
			FromUid:           listenerUids[idx],
			ToUid:             uid,
			MsgType:           notify.TextMsgTypeHelloMsg2,
			Title:             "",
			Val1:              "",
			Val2:              "",
			Val3:              "",
			Val4:              "",
			Val5:              "",
			Val6:              "",
			Sync:              tim.TimMsgSyncFromNo,
			RepeatMsgCheckId:  "",
			RepeatMsgCheckSec: 0,
		}

		kqm2.Text = l.GetHelloMsg(uid, listenerUids[idx], hm[idx])
		if kqm2.Text == "" {
			continue
		}
		msgs[idx] = &kqm2
	}

	for idx := 0; idx < len(msgs); idx++ {
		if msgs[idx] == nil {
			continue
		}
		var buf []byte
		buf, err = json.Marshal(msgs[idx])
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic json Marshal err:%+v", err)
			return
		}
		var payload []byte
		payload, err = json.Marshal(jobtype.DeferSendImMsgPayload{KqMsgBuf: buf})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic json Marshal err:%+v", err)
			return
		}
		delaySecond := notify.DeferSendNewUserRecommendListenerImMsgSecond*(int64(idx)+1) + (uid % notify.DeferSendNewUserRecommendListenerImMsgSecond)
		_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferSendImMsg, payload), asynq.ProcessIn(time.Duration(delaySecond)*time.Second))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic EnqueueContext err:%+v", err)
			return
		}
	}

	// 设置消息发送标记
	_, err = l.svcCtx.ImRedis.SetNotifyLimit(l.ctx, fmt.Sprintf("%d-%d", uid, notify.TextMsgTypeHelloMsg2), notify.TextMsgTypeHelloMsg2, notify.SendListenerHelloMsgIntervalSecond)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic SetNotifyLimit err:%+v", err)
		return
	}
	return
}

// 非第一次登陆
func (l *SendHelloWhenUserLoginMq) sendUserMsg(uid int64, listenerUids []int64) {
	if len(listenerUids) <= 0 {
		return
	}

	// 获取聊天消息
	var hlm []int64
	var err error
	hlm, err = l.svcCtx.ListenerRedis.PopUserHelloMsg(l.ctx, uid, int64(len(listenerUids)))
	if err != nil {
		return
	}
	if len(hlm) <= 0 {
		return
	}

	msgs := make([]*kqueue.SendImDefineMessage, len(listenerUids))
	for idx := 0; idx < len(listenerUids); idx++ {
		kqm := kqueue.SendImDefineMessage{
			FromUid:           listenerUids[idx],
			ToUid:             uid,
			MsgType:           notify.TextMsgTypeHelloMsg2,
			Title:             "",
			Val1:              "",
			Val2:              "",
			Val3:              "",
			Val4:              "",
			Val5:              "",
			Val6:              "",
			Sync:              tim.TimMsgSyncFromNo,
			RepeatMsgCheckId:  "",
			RepeatMsgCheckSec: 0,
		}

		kqm.Text = l.GetHelloMsg(uid, listenerUids[idx], hlm[idx])
		if kqm.Text == "" {
			continue
		}
		msgs[idx] = &kqm
	}

	for idx := 0; idx < len(msgs); idx++ {
		if msgs[idx] == nil {
			continue
		}
		var buf []byte
		buf, err = json.Marshal(msgs[idx])
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic sendUserMsg json Marshal err:%+v", err)
			return
		}
		var payload []byte
		payload, err = json.Marshal(jobtype.DeferSendImMsgPayload{KqMsgBuf: buf})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic sendUserMsg json Marshal err:%+v", err)
			return
		}
		delaySecond := notify.DeferSendUserRecommendListenerImMsgSecond*(int64(idx)+1) + (uid % notify.DeferSendUserRecommendListenerImMsgIntervalSecond) + notify.DeferSendUserRecommendListenerImMsgIntervalSecond
		_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferSendImMsg, payload), asynq.ProcessIn(time.Duration(delaySecond)*time.Second))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic sendUserMsg EnqueueContext err:%+v", err)
			return
		}
	}

	// 设置消息发送标记
	_, err = l.svcCtx.ImRedis.SetNotifyLimit(l.ctx, fmt.Sprintf("%d-%d", uid, notify.TextMsgTypeHelloMsg2), notify.TextMsgTypeHelloMsg2, notify.SendListenerHelloMsgIntervalSecond)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic sendUserMsg SetNotifyLimit err:%+v", err)
		return
	}
	return
}

func (l *SendHelloWhenUserLoginMq) GetHelloMsg(uid, listenerUid, idx int64) string {
	if int(idx) >= len(notify.LoginHelloMsg2) {
		return ""
	}

	tmp := notify.LoginHelloMsg2[idx]

	switch idx {
	case 1: // 用户名称
		rsp1, err := l.svcCtx.UsercenterRpc.GetUserShortProfile(l.ctx, &pbUser.GetUserShortProfileReq{Uid: uid})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetHelloMsg GetUserShortProfile err:%+v", err)
			return ""
		}
		return fmt.Sprintf(tmp, rsp1.User.Nickname)

	case 2, 5, 7, 9, 11: // XXX名称
		rsp1, err := l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &pbListener.GetListenerBasicInfoReq{ListenerUid: listenerUid})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetHelloMsg GetListenerBasicInfo err:%+v", err)
			return ""
		}
		nm := strings.ReplaceAll(rsp1.NickName, "#", "")
		return fmt.Sprintf(tmp, nm)
	default:
		return tmp
	}
}
