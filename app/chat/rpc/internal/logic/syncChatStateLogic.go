package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/chatPgModel"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"strconv"
	"time"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncChatStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncChatStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncChatStateLogic {
	return &SyncChatStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  同步聊天状态过程
func (l *SyncChatStateLogic) SyncChatState(in *pb.SyncChatStateReq) (resp *pb.SyncChatStateResp, err error) {
	resp = &pb.SyncChatStateResp{}
	if in.Uid == in.ListenerUid {
		return
	}

	switch in.Action {
	case chatkey.ChatAction1:
		resp, err = l.enterChat(in)
	case chatkey.ChatAction2:
		resp, err = l.updateTextChatFreeCnt(in)
	case chatkey.ChatAction3, chatkey.ChatAction7:
		resp, err = l.startVoiceChat(in)
	case chatkey.ChatAction11:
		resp, err = l.check(in)
	case chatkey.ChatAction12:
		resp, err = l.queryBalance(in)
	case chatkey.ChatAction4, chatkey.ChatAction5, chatkey.ChatAction6, chatkey.ChatAction8, chatkey.ChatAction9, chatkey.ChatAction10:
		resp, err = l.stopVoiceChat(in)
	case chatkey.ChatAction13:
		resp, err = l.startCallListener(in)
	case chatkey.ChatAction14:
		resp, err = l.releaseLockListener(in)

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "wrong action")
	}
	if err != nil {
		return nil, err
	}
	if in.Action != chatkey.ChatAction2 {
		var b bool
		resp.FreeChatCnt, err, b = l.svcCtx.ChatRedis.GetSetListenerFreeMsgCnt(l.ctx, in.Uid, in.ListenerUid)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic GetSetListenerFreeMsgCnt err:%+v", err)
			return nil, err
		}
		if b {
			err = l.svcCtx.ChatBalanceModel.UpdateEnterChatTime(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
			if err != nil {
				return nil, err
			}
		}
	}

	if resp.FreeChatCnt < 0 {
		resp.FreeChatCnt = 0
	}
	// 获取当前XXX的通话状态
	if resp.ListenerChatState == 0 {
		var lvc *chatPgModel.ListenerVoiceChatState
		lvc, err = l.svcCtx.ListenerVoiceChatStateModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic ListenerVoiceChatStateModel.FindOne err:%+v", err)
			return nil, err
		}
		resp.ListenerChatState = lvc.State
	}

	// 判断当前用户状态
	var textChatExTime time.Time // 文字聊天到期时间
	if resp.TextChatExpiryTime != "" {
		textChatExTime, err = time.ParseInLocation(db.DateTimeFormat, resp.TextChatExpiryTime, time.Local)
		if err != nil {
			return nil, err
		}
	}

	if resp.VoiceChatMinute <= 0 && resp.TextChatExpiryTime == "" && resp.UsedVoiceChatMinute <= 0 {
		resp.ChatState = chatkey.UserChatState1
	} else if resp.VoiceChatMinute > 0 || (resp.TextChatExpiryTime != "" && textChatExTime.After(time.Now())) {
		resp.ChatState = chatkey.UserChatState2
	} else if resp.UsedVoiceChatMinute > 0 || resp.TextChatExpiryTime != "" {
		resp.ChatState = chatkey.UserChatState3
	}

	if textChatExTime.Before(time.Now()) { // 处理文字聊天到期时间
		resp.TextChatExpiryTime = ""
	}

	resp.Uid = in.Uid
	resp.ListenerUid = in.ListenerUid
	return resp, nil
}

func (l *SyncChatStateLogic) startCallListener(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	rsp, err := l.getBalance(in)
	if err != nil {
		return nil, err
	}
	if rsp.VoiceChatMinute > 0 && tool.IsInt64ArrayExist(rsp.CurrentVoiceChatState, []int64{0, chatkey.VoiceChatStateSettle}) {
		var lvc *chatPgModel.ListenerVoiceChatState
		lvc, err = l.svcCtx.ListenerVoiceChatStateModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic ListenerVoiceChatStateModel.FindOne err:%+v", err)
			return nil, err
		}
		rsp.ListenerChatState = lvc.State
		if tool.IsInt64ArrayExist(rsp.ListenerChatState, []int64{0, chatkey.VoiceChatStateSettle}) {
			err = l.lockListener(in.Uid, in.ListenerUid)
			if err != nil {
				return nil, err
			}
		}
	}

	return rsp, err
}

func (l *SyncChatStateLogic) releaseLockListener(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	rsp, err := l.getBalance(in)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ChatRedis.ReleaseLockTargetUser(l.ctx, in.Uid, in.ListenerUid)
	if err != nil {
		return nil, err
	}
	return rsp, err
}

func (l *SyncChatStateLogic) lockListener(fromUid, toUid int64) error {
	if toUid == 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "no listener uid")
	}
	// 锁定
	b, err := l.svcCtx.ChatRedis.RecordLockTargetUser(l.ctx, fromUid, toUid, chatkey.StartCallLockTargetSeconds)
	if err != nil {
		return err
	}
	if !b {
		return xerr.NewGrpcErrCodeMsg(xerr.ChatErrorListenerBusy, "对方占线")
	}
	return nil
}

func (l *SyncChatStateLogic) enterChat(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}
	rs, isFirst, err := l.getChatBalance(in)
	if err != nil {
		return nil, err
	}

	resp := getChatBalanceResp(rs)
	resp.IsFirstEnterChat = isFirst

	if rs.CurrentVoiceChatState == chatkey.VoiceChatStateStop {
		l.checkVoiceChat(rs, 3)
	}
	return resp, nil
}

func (l *SyncChatStateLogic) check(in *pb.SyncChatStateReq) (resp *pb.SyncChatStateResp, err error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}
	var data *chatPgModel.ChatBalance
	var isFirst int64
	data, isFirst, err = l.getChatBalance(in)
	if err != nil {
		return nil, err
	}
	resp = getChatBalanceResp(data)
	resp.IsFirstEnterChat = isFirst
	l.checkVoiceChat(data, 0)
	return
}

func getChatBalanceResp(data *chatPgModel.ChatBalance) (resp *pb.SyncChatStateResp) {
	rsp := pb.SyncChatStateResp{
		FreeChatCnt:           0,
		TextChatExpiryTime:    "",
		VoiceChatMinute:       0,
		UsedVoiceChatMinute:   0,
		ChatState:             0,
		ListenerChatState:     0,
		Uid:                   0,
		ListenerUid:           0,
		CurrentVoiceChatState: 0,
		IsFirstEnterChat:      0,
	}
	if data != nil {
		if data.TextChatExpiryTime.Valid {
			rsp.TextChatExpiryTime = data.TextChatExpiryTime.Time.Format(db.DateTimeFormat)
		}
		rsp.VoiceChatMinute = data.AvailableVoiceChatMinute
		rsp.UsedVoiceChatMinute = data.UsedVoiceChatMinute
		rsp.CurrentVoiceChatState = data.CurrentVoiceChatState
		return &rsp
	}
	return
}

func (l *SyncChatStateLogic) getBalance(in *pb.SyncChatStateReq) (resp *pb.SyncChatStateResp, err error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}
	rs, isFirst, err := l.getChatBalance(in)
	if err != nil {
		return nil, err
	}
	resp = getChatBalanceResp(rs)
	resp.IsFirstEnterChat = isFirst
	return
}

func (l *SyncChatStateLogic) queryBalance(in *pb.SyncChatStateReq) (resp *pb.SyncChatStateResp, err error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}
	var data *chatPgModel.ChatBalance
	data, err = l.svcCtx.ChatBalanceModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
	if err != nil && err != chatPgModel.ErrNotFound {
		return
	}
	err = nil
	if data != nil {
		resp = getChatBalanceResp(data)
	} else { // 第一次
		resp = &pb.SyncChatStateResp{
			FreeChatCnt:           0,
			TextChatExpiryTime:    "",
			VoiceChatMinute:       0,
			UsedVoiceChatMinute:   0,
			ChatState:             0,
			ListenerChatState:     0,
			Uid:                   0,
			ListenerUid:           0,
			CurrentVoiceChatState: 0,
			IsFirstEnterChat:      0,
		}
	}

	return
}

func (l *SyncChatStateLogic) getChatBalance(in *pb.SyncChatStateReq) (data *chatPgModel.ChatBalance, isFirstEnterChat int64, err error) {
	data, err = l.svcCtx.ChatBalanceModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
	if err != nil && err != chatPgModel.ErrNotFound {
		return
	}
	if data == nil && err == chatPgModel.ErrNotFound {
		// 新建 chat_balance
		data = new(chatPgModel.ChatBalance)
		data.Uid = in.Uid
		data.ListenerUid = in.ListenerUid
		data.Id = fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid)
		_, err = l.svcCtx.ChatBalanceModel.Insert(l.ctx, data)
		if err != nil {
			return
		}

		// 创建 user_listener_relation
		ulr := chatPgModel.UserListenerRelation{
			Id:          fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid),
			Uid:         in.Uid,
			ListenerUid: in.ListenerUid,
			TotalScore:  0,
		}
		_, err = l.svcCtx.UserListenerRelationModel.Insert(l.ctx, &ulr)
		if err != nil {
			return
		}

		// 第一次点进聊天 发送个人介绍
		isFirstEnterChat = db.Enable
		//
		l.userEnterChat(in, isFirstEnterChat)
		return
	}
	// 非首次进入聊天
	if in.Action == chatkey.ChatAction1 {
		l.userEnterChat(in, 0)
	}
	return
}

func (l *SyncChatStateLogic) userEnterChat(in *pb.SyncChatStateReq, isFirst int64) {
	kqm := kqueue.UserFirstEnterChatMessage{
		Uid:         in.Uid,
		ListenerUid: in.ListenerUid,
		IsFirst:     isFirst,
	}
	buf, err := json.Marshal(&kqm)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SyncChatStateLogic userEnterChat json Marshal err:%+v", err)
		return
	}

	err = l.svcCtx.KqueueFirstEnterChatClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SyncChatStateLogic userEnterChat kafka Push err:%+v", err)
		return
	}
	return
}

// 检查异常通话
func (l *SyncChatStateLogic) checkVoiceChat(data *chatPgModel.ChatBalance, delayMin int64) {
	if data.CurrentVoiceChatState == chatkey.VoiceChatStateSettle { // 正常状态
		return
	}
	// 检查用户和XXX的通话状态 未结算状态
	var voiceExpiryTime *time.Time
	if data.CurrentVoiceChatExpiryTime.Valid {
		voiceExpiryTime = &(data.CurrentVoiceChatExpiryTime.Time)
	}
	// 上次结束时间
	var stopTime *time.Time
	if data.CurrentStopTime.Valid {
		stopTime = &(data.CurrentStopTime.Time)
	}
	bTime := time.Now()
	if delayMin > 0 {
		bTime = bTime.Add(-time.Duration(delayMin) * time.Minute)
	}

	if data.CurrentVoiceChatState == chatkey.VoiceChatStateStart && voiceExpiryTime.Before(bTime) { // 未停止
		_, err := l.updateChatStop(data.Uid, data.ListenerUid, data.CurrentChatLogId, chatkey.ChatAction11, chatkey.VoiceChatStateStop, &(data.CurrentStartTime.Time), &(data.CurrentVoiceChatExpiryTime.Time))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("checkVoiceChat updateChatStop error chat log id %d, update stop error:%+v", data.CurrentChatLogId, err)
			return
		}
	} else if data.CurrentVoiceChatState == chatkey.VoiceChatStateStop && stopTime.Before(bTime) { // 未结算
		l.settle(data.Uid, data.ListenerUid, data.CurrentChatLogId, &(data.CurrentStartTime.Time), &(data.CurrentStopTime.Time))
	}
	return
}

// 文字聊天次数减1
func (l *SyncChatStateLogic) updateTextChatFreeCnt(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}

	fc, err := l.svcCtx.ChatRedis.DecrListenerFreeMsgCnt(l.ctx, in.Uid, in.ListenerUid)
	if err != nil {
		return nil, err
	}

	var rsp *pb.SyncChatStateResp
	rsp, err = l.getBalance(in)
	if err != nil {
		return nil, err
	}
	rsp.FreeChatCnt = fc
	return rsp, nil
}

func (l *SyncChatStateLogic) startVoiceChat(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数为空 %d %d", in.Uid, in.ListenerUid))
	}
	rs1, err := l.svcCtx.ChatBalanceModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
	if err != nil {
		return nil, err
	}
	state := chatkey.VoiceChatStateStart
	var voiceExpiryTime *time.Time
	if rs1.CurrentVoiceChatExpiryTime.Valid {
		voiceExpiryTime = &rs1.CurrentVoiceChatExpiryTime.Time
	}
	if rs1.CurrentVoiceChatState == chatkey.VoiceChatStateStart { // 还有语音时间 已经开始
		if voiceExpiryTime.After(time.Now()) {
			resp := &pb.SyncChatStateResp{
				VoiceChatMinute: rs1.AvailableVoiceChatMinute,
			}
			if rs1.TextChatExpiryTime.Valid {
				resp.TextChatExpiryTime = rs1.TextChatExpiryTime.Time.Format(db.DateTimeFormat)
			}
			resp.UsedVoiceChatMinute = rs1.UsedVoiceChatMinute
			resp.ListenerChatState = state
			resp.CurrentVoiceChatState = rs1.CurrentVoiceChatState
			return resp, nil
		}
		// 语音通话已经用尽 但是状态还是开始状态 异常状态需要处理
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("error chat state %d", rs1.CurrentVoiceChatState))
	}

	resp := &pb.SyncChatStateResp{}
	resp.ListenerChatState = state
	resp.CurrentVoiceChatState = state
	err = l.svcCtx.ChatBalanceModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 通话log
		data := new(chatPgModel.VoiceChatLog)
		data.ListenerUid = in.ListenerUid
		data.Uid = in.Uid
		data.StartTime = time.Now()
		data.StartAction = in.Action
		data.State = state
		data.Id = uniqueid.GenDataId()
		err2 := l.svcCtx.VoiceChatLogModel.Insert2(ctx, session, data)
		if err2 != nil {
			return err2
		}
		// 更新通话状态
		err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, resp.UsedVoiceChatMinute = l.svcCtx.ChatBalanceModel.UpdateCurrentChatStart(ctx, session, in.Uid, in.ListenerUid, data.Id, &data.StartTime)
		if err2 != nil {
			return err2
		}
		// 分别更新用户和XXX当前语音通话状态
		err2 = l.svcCtx.UserVoiceChatStateModel.UpdateState(ctx, session, in.Uid, in.ListenerUid, state)
		if err2 != nil {
			return err2
		}
		err2 = l.svcCtx.ListenerVoiceChatStateModel.UpdateState(ctx, session, in.Uid, in.ListenerUid, state)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发送im消息
	kqMsg := kqueue.SendImDefineMessage{
		FromUid:           in.Uid,
		ToUid:             in.ListenerUid,
		MsgType:           notify.DefineNotifyMsgTypeChatMsg6,
		Title:             "",
		Text:              notify.DefineNotifyMsgTemplateChatMsg6,
		Val1:              strconv.FormatInt(in.Uid, 10),
		Val2:              strconv.FormatInt(in.ListenerUid, 10),
		Val3:              "",
		Val4:              "",
		Val5:              "",
		Val6:              "",
		Sync:              tim.TimMsgSyncFromYes,
		RepeatMsgCheckId:  "",
		RepeatMsgCheckSec: 0,
	}

	buf, err := json.Marshal(&kqMsg)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *SyncChatStateLogic) stopVoiceChat(in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	// 由于客户端技术问题 有可能无法获取到成对的uid和listenerUid 需要从当前通话状态中查询
	if in.Uid != 0 && in.ListenerUid == 0 {
		rs, err := l.svcCtx.UserVoiceChatStateModel.FindOne(l.ctx, in.Uid)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic stopVoiceChat UserVoiceChatStateModel FindOne in:%+v err:%+v", in, err)
			return nil, err
		}
		if rs.ListenerUid == 0 {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic stopVoiceChat ListenerUid is empty in:%+v err:%+v", in, err)
			return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, "返回空缓存")
		}
		in.ListenerUid = rs.ListenerUid
	} else if in.Uid == 0 && in.ListenerUid != 0 {
		rs, err := l.svcCtx.ListenerVoiceChatStateModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic stopVoiceChat ListenerVoiceChatStateModel FindOne in:%+v err:%+v", in, err)
			return nil, err
		}
		if rs.Uid == 0 {
			logx.WithContext(l.ctx).Errorf("SyncChatStateLogic stopVoiceChat Uid is empty in:%+v err:%+v", in, err)
			return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, "返回空缓存")
		}
		in.Uid = rs.Uid
	} else if in.Uid == 0 && in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "uid both empty")
	}
	//
	rs, err := l.svcCtx.ChatBalanceModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
	if err != nil {
		return nil, err
	}
	state := chatkey.VoiceChatStateStop
	if tool.IsInt64ArrayExist(rs.CurrentVoiceChatState, []int64{chatkey.VoiceChatStateStop, chatkey.VoiceChatStateSettle}) { // 已经结束
		resp := &pb.SyncChatStateResp{
			VoiceChatMinute: rs.AvailableVoiceChatMinute,
		}
		if rs.TextChatExpiryTime.Valid {
			resp.TextChatExpiryTime = rs.TextChatExpiryTime.Time.Format(db.DateTimeFormat)
		}
		resp.ListenerChatState = state
		resp.UsedVoiceChatMinute = rs.UsedVoiceChatMinute
		resp.CurrentVoiceChatState = rs.CurrentVoiceChatState
		return resp, nil
	}
	now := time.Now()
	return l.updateChatStop(rs.Uid, rs.ListenerUid, rs.CurrentChatLogId, in.Action, state, &(rs.CurrentStartTime.Time), &now)
}

func (l *SyncChatStateLogic) updateChatStop(uid, listenerUid int64, chatLogId string, action, state int64, startTime, stopTime *time.Time) (*pb.SyncChatStateResp, error) {
	resp := &pb.SyncChatStateResp{}
	resp.ListenerChatState = state
	resp.CurrentVoiceChatState = state
	err := l.svcCtx.ChatBalanceModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 分别更新用户和XXX当前语音通话状态
		err2 := l.svcCtx.UserVoiceChatStateModel.UpdateState(ctx, session, uid, listenerUid, state)
		if err2 != nil {
			return err2
		}
		err2 = l.svcCtx.ListenerVoiceChatStateModel.UpdateState(ctx, session, uid, listenerUid, state)
		if err2 != nil {
			return err2
		}
		err2 = l.svcCtx.VoiceChatLogModel.UpdateChatLog(ctx, session, chatLogId, action, state, stopTime)
		if err2 != nil {
			return err2
		}
		err2, resp.TextChatExpiryTime, resp.VoiceChatMinute, resp.UsedVoiceChatMinute = l.svcCtx.ChatBalanceModel.UpdateCurrentChatStop(ctx, session, uid, listenerUid, chatLogId, stopTime)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 结算
	l.settle(uid, listenerUid, chatLogId, startTime, stopTime)
	return resp, nil
}

func (l *SyncChatStateLogic) settle(uid, listenerUid int64, chatLogId string, startTime, stopTime *time.Time) {
	// 更新本次通话的使用情况
	var msg kqueue.UpdateChatStatMessage
	msg.LogId = chatLogId
	msg.Uid = uid
	msg.ListenerUid = listenerUid
	msg.StartTime = startTime.Format(db.DateTimeFormat)
	msg.StopTime = stopTime.Format(db.DateTimeFormat)
	msg.OrderType = orderkey.ListenerOrderTypeVoiceChat

	// 校验时间
	if stopTime.Unix()-startTime.Unix() <= 0 {
		return
	}

	buf, err := json.Marshal(&msg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SyncChatStateLogic settle json marshal chat log id:%d err:%+v", chatLogId, err)
		return
	}
	err = l.svcCtx.KqueueUpdateChatStatClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SyncChatStateLogic settle KqueueUpdateChatStatClient push chat log id:%d err:%+v", chatLogId, err)
		return
	}
	return
}
