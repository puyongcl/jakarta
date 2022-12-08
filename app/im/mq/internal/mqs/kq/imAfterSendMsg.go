package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	"jakarta/app/im/mq/internal/svc"
	pbIm "jakarta/app/im/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/userkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"strconv"
	"time"
)

// 发送自定义消息
type ImAfterSendMsgMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImAfterSendMsgMq(ctx context.Context, svcCtx *svc.ServiceContext) *ImAfterSendMsgMq {
	return &ImAfterSendMsgMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImAfterSendMsgMq) Consume(_, val string) error {
	var message kqueue.ImAfterSendMsg
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("ImAfterSendMsgMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("ImAfterSendMsgMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *ImAfterSendMsgMq) execService(message *kqueue.ImAfterSendMsg) error {
	// 如果是非用户账号 则不记录
	fromUid, err := strconv.ParseInt(message.FromUid, 10, 64)
	if err != nil {
		return err
	}
	toUid, err := strconv.ParseInt(message.ToUid, 10, 64)
	if err != nil {
		return err
	}
	if fromUid < userkey.UidStart || toUid < userkey.UidStart {
		return nil
	}
	// 获取用户类型
	fromUserType, toUserType, err := l.GetUserType(fromUid, toUid)
	if err != nil {
		return err
	}
	if fromUserType == toUserType {
		logx.WithContext(l.ctx).Errorf("error user type %d %d", fromUid, toUid)
		return nil
	}

	et := pbChat.SendUserListenerRelationEventReq{}

	if fromUserType == userkey.UserTypeListener { // 发送者是XXX
		// 发送方是XXX 回复消息 移除标记
		err = l.fromUserIsListener(message.FromUid)
		if err != nil {
			return err
		}

		// 发送者是XXX 延长 浏览通知发送时间
		err = l.svcCtx.IMRedis.RenewNotifyLimit(l.ctx, fmt.Sprintf("%s-%s", message.ToUid, message.FromUid), notify.DefineNotifyMsgTypeViewMsg17, notify.SendViewNotifyLimitRenewMin*60)
		if err != nil {
			return err
		}

		et.Uid = toUid
		et.ListenerUid = fromUid
		et.EventType = chatkey.InteractiveEventTypeSendMsg2
	}
	if toUserType == userkey.UserTypeListener { // 接收者是XXX
		// 接收方是XXX 记录时间 并判断是否已经标记 如果超时 则转状态为 休息中
		err = l.toUserIsListener(toUid)
		if err != nil {
			return err
		}

		// 接收者是XXX 延长 浏览通知发送时间
		err = l.svcCtx.IMRedis.RenewNotifyLimit(l.ctx, fmt.Sprintf("%s-%s", message.FromUid, message.ToUid), notify.DefineNotifyMsgTypeViewMsg17, notify.SendViewNotifyLimitRenewMin*60)
		if err != nil {
			return err
		}

		et.Uid = fromUid
		et.ListenerUid = toUid
		et.EventType = chatkey.InteractiveEventTypeSendMsg1
	}

	// 发送小程序订阅通知
	var content string
	switch message.MsgType {
	case tim.TIMSoundElem:
		content = notify.DefineNotifyMsgTypeMiniProgramMsg1VoiceMsgDefaultContent
	case tim.TIMImageElem:
		content = notify.DefineNotifyMsgTypeMiniProgramMsg1ImageMsgDefaultContent
	default:
		content = message.Text
	}

	err = l.mpNotifyNewMsg(fromUid, toUid, fromUserType, content)
	if err != nil {
		return err
	}

	// log
	in := pbIm.AddImMsgLogReq{
		FromUid:      fromUid,
		ToUid:        toUid,
		MsgTime:      message.MsgTime,
		MsgId:        message.MsgKey,
		MsgType:      message.MsgType,
		MsgSeq:       message.MsgSeq,
		FromUserType: fromUserType,
	}
	_, _ = l.svcCtx.ImRpc.AddImMsgLog(l.ctx, &in)

	// 发送用户和XXX互动事件
	_, _ = l.svcCtx.ChatRpc.SendUserListenerRelationEvent(l.ctx, &et)
	return nil
}

// 当用户有订阅小程序通知 不在线 发送留言通知
func (l *ImAfterSendMsgMq) mpNotifyNewMsg(fromUid, toUid, fromUserType int64, content string) error {
	// 判断用户是否有订阅
	b, err := l.svcCtx.IMRedis.IsUserHaveOneTimeSubscribeMsg(l.ctx, notify.DefineNotifyMsgTypeMiniProgramMsg1, toUid)
	if err != nil {
		return err
	}
	if !b {
		return nil
	}
	// 判断用户是否在线
	//var rsp *pbUser.GetUserOnlineStateResp
	//rsp, err = l.svcCtx.UserRpc.GetUserOnlineState(l.ctx, &pbUser.GetUserOnlineStateReq{Uid: toUid})
	//if err != nil {
	//	return err
	//}
	//if rsp.OnlineState == userkey.Login { // 在线
	//	return nil
	//}
	// 获取用户资料
	var nickname string
	if fromUserType == userkey.UserTypeListener {
		var rsp2 *pbListener.GetListenerBasicInfoResp
		rsp2, err = l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &pbListener.GetListenerBasicInfoReq{
			ListenerUid: fromUid,
		})
		if err != nil {
			return err
		}
		nickname = rsp2.NickName
	} else {
		var rsp2 *pbUser.GetUserShortProfileResp
		rsp2, err = l.svcCtx.UserRpc.GetUserShortProfile(l.ctx, &pbUser.GetUserShortProfileReq{Uid: fromUid})
		if err != nil {
			return err
		}
		if rsp2.User != nil {
			nickname = rsp2.User.Nickname
		}
	}

	// 发送消息
	mpMsg := kqueue.SendMiniProgramSubscribeMessage{
		Thing4:  content,
		Thing5:  nickname,
		Time3:   time.Now().Format(db.DateTimeFormat),
		MsgType: notify.DefineNotifyMsgTypeMiniProgramMsg1,
		ToUid:   toUid,
		Page:    notify.DefineNotifyMsgTypeMiniProgramMsg1Path,
	}
	kqMsg := kqueue.SubscribeNotifyMsgMessage{
		Uid:       toUid,
		TargetUid: 0,
		MsgType:   notify.DefineNotifyMsgTypeMiniProgramMsg1,
		SendCnt:   0,
		Action:    notify.SubscribeOneTimeNotifyMsgEventSend,
		IMMsg:     nil,
		MpMsg:     &mpMsg,
	}
	var buf []byte
	buf, err = json.Marshal(kqMsg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueSendSubscribeNotifyMsgClient.Push(string(buf))
	if err != nil {
		return err
	}
	return nil
}

// 发送方是XXX 回复消息 移除标记
func (l *ImAfterSendMsgMq) fromUserIsListener(fromUid string) error {
	return l.svcCtx.IMRedis.CancelMarkUserReplyTime(l.ctx, fromUid)
}

// 接收方是XXX 记录时间 并判断是否已经标记 如果超时 则转状态为 休息中
func (l *ImAfterSendMsgMq) toUserIsListener(toUid int64) error {
	//state, onlineState, err := l.GetListenerWorkState(toUid)
	state, _, err := l.GetListenerWorkState(toUid)
	if err != nil {
		return err
	}
	// 已经在休息中
	if tool.IsInt64ArrayExist(state, listenerkey.ListenerRestState) {
		return nil
	}
	now := time.Now()
	// 加入标记
	stamp, err := l.svcCtx.IMRedis.MarkUserReplyTime(l.ctx, toUid, now.Format(db.DateTimeFormat))
	if err != nil {
		return err
	}
	// 判断是否标记过 是否过期
	if stamp == "" { // 首次标记
		return nil
	}
	t, err := time.ParseInLocation(db.DateTimeFormat, stamp, time.Local)
	if err != nil {
		return err
	}
	if now.Sub(t).Minutes() < listenerkey.NotReplyAutoSwitchWorkStateIntervalMinute {
		return nil
	}
	// 超过4分钟未回复 如果XXX不在线 转为休息中
	//// 判断XXX是否在线
	//if onlineState == userkey.Login { // 在线
	//	return nil
	//}
	//// 不在线 则转状态为 休息中
	_, err = l.svcCtx.ListenerRpc.ChangeWorkState(l.ctx, &pbListener.ChangeWorkStateReq{
		ListenerUid: toUid,
		WorkState:   listenerkey.ListenerWorkStateRestingAuto,
	})
	if err != nil {
		return err
	}
	// 移除标记
	err = l.svcCtx.IMRedis.CancelMarkUserReplyTime(l.ctx, strconv.FormatInt(toUid, 10))
	if err != nil {
		return err
	}

	// 服务号通知XXX
	l.fwhNotifyListener(toUid)
	return nil
}

func (l *ImAfterSendMsgMq) fwhNotifyListener(toUid int64) {
	fwhMsg := kqueue.SendFwhNotifyMessage{
		First:    fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg5FirstData, listenerkey.NotReplyAutoSwitchWorkStateIntervalMinute),
		Keyword1: notify.DefineNotifyMsgTypeFwhMsg5Key1,
		Keyword2: notify.DefineNotifyMsgTypeFwhMsg5Key2,
		Keyword3: time.Now().Format(db.DateTimeFormat),
		Remark:   notify.DefineNotifyMsgTypeFwhMsg5RemarkData,
		MsgType:  notify.DefineNotifyMsgTypeFwhMsg5,
		Color:    notify.DefineNotifyMsgTypeFwhMsg5Color,
		ToUid:    toUid,
	}

	buf, err := json.Marshal(&fwhMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ImAfterSendMsgMq notify json marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueSendWxFwhMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ImAfterSendMsgMq notify kafka Push err:%+v", err)
		return
	}
}

func (l *ImAfterSendMsgMq) GetUserType(fromUid, toUid int64) (fromUserType int64, toUserType int64, err error) {
	r1, err := l.svcCtx.ImMemoryCache.GetUserInfo(l.ctx, fromUid)
	if err != nil {
		return
	}
	r2, err := l.svcCtx.ImMemoryCache.GetUserInfo(l.ctx, toUid)
	if err != nil {
		return
	}
	return r1.UserType, r2.UserType, nil
}

func (l *ImAfterSendMsgMq) GetListenerWorkState(uid int64) (int64, int64, error) {
	rsp, err := l.svcCtx.ListenerRpc.GetWorkState(l.ctx, &pbListener.GetWorkStateReq{ListenerUid: uid})
	if err != nil {
		return 0, 0, err
	}
	return rsp.WorkState, rsp.OnlineState, nil
}
