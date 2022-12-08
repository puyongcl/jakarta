package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeWorkStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeWorkStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeWorkStateLogic {
	return &ChangeWorkStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改工作状态和休息时间
func (l *ChangeWorkStateLogic) ChangeWorkState(in *pb.ChangeWorkStateReq) (*pb.ChangeWorkStateResp, error) {
	data := new(listenerPgModel.ListenerProfile)
	data.ListenerUid = in.ListenerUid
	if in.WorkState != 0 {
		data.WorkState = in.WorkState
	}
	if in.RestingTimeEnable != 0 {
		data.RestingTimeEnable = in.RestingTimeEnable
	}
	if in.StartWorkTime != "" {
		data.StartWorkTime = in.StartWorkTime
	}
	if in.StopWorkTime != "" {
		data.StopWorkTime = in.StopWorkTime
	}
	if len(in.WorkDays) >= 0 {
		data.WorkDays = in.WorkDays
	}
	err := l.svcCtx.ListenerProfileModel.UpdateWorkState(l.ctx, nil, data)
	if err != nil {
		return nil, err
	}

	// 通知订阅消息 暂不做
	switch in.WorkState {
	case listenerkey.ListenerWorkStateWorking: // 可接单 通知订阅用户 TODO 暂不做
	case listenerkey.ListenerWorkStateRestingAuto: // 自动转为休息中 系统通知XXX
		l.notify2(in)
	default:

	}

	// 修改新用户推荐列表
	switch in.WorkState {
	case listenerkey.ListenerWorkStateWorking: // 改为工作中
		err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendWhenUserLogin, in.ListenerUid)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendReplyUserStory, in.ListenerUid)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendSendHelloMsgWhenUserLogin, in.ListenerUid)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerWorkStateRestingAuto, listenerkey.ListenerWorkStateRestingManual, listenerkey.ListenerWorkStateAccountDeleted: // 休息中、删号
		err = l.svcCtx.ListenerRedis.RemoveNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendWhenUserLogin, in.ListenerUid)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerRedis.RemoveNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendReplyUserStory, in.ListenerUid)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerRedis.RemoveNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendSendHelloMsgWhenUserLogin, in.ListenerUid)
		if err != nil {
			return nil, err
		}

	default:

	}

	return &pb.ChangeWorkStateResp{}, nil
}

// 通知XXX已经自动转为休息中
func (l *ChangeWorkStateLogic) notify2(in *pb.ChangeWorkStateReq) {
	kqMsg := kqueue.SendImDefineMessage{
		FromUid:           notify.TimSystemNotifyUid,
		ToUid:             in.ListenerUid,
		MsgType:           notify.DefineNotifyMsgTypeSystemMsg27,
		Title:             notify.DefineNotifyMsgTemplateSystemMsgTitle27,
		Text:              notify.DefineNotifyMsgTemplateSystemMsg27,
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
	buf, err := json.Marshal(kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ChangeWorkStateLogic notify2 json Marshal uid:%d err:%+v", in.ListenerUid, err)
		return
	}

	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ChangeWorkStateLogic notify2 kafka Push uid:%d err:%+v", in.ListenerUid, err)
		return
	}
}

// 通知订阅可接单用户
func (l *ChangeWorkStateLogic) notify(in *pb.ChangeWorkStateReq) {
	data, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ChangeWorkStateLogic notify ListenerProfileModel FindOne uid:%d err:%+v", in.ListenerUid, err)
		return
	}
	kqSendImMsg := kqueue.SendImDefineMessage{
		FromUid: notify.TimSystemNotifyUid,
		ToUid:   0,
		MsgType: notify.DefineNotifyMsgTypeSystemMsg23,
		Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle23,
		Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsg23, strings.Replace(data.NickName, "#", "", -1)),
		Val1:    strconv.FormatInt(in.ListenerUid, 10),
		Val2:    "",
		Sync:    tim.TimMsgSyncFromNo,
	}

	kqNotifyMsg := kqueue.SubscribeNotifyMsgMessage{
		Uid:       0,
		TargetUid: in.ListenerUid,
		MsgType:   notify.DefineNotifyMsgTypeSystemMsg23,
		SendCnt:   notify.SubscribeNotifyMsgSendCntOne,
		Action:    notify.SubscribeUserNotifyMsgEventSend,
		IMMsg:     &kqSendImMsg,
	}
	var buf []byte
	buf, err = json.Marshal(kqNotifyMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ChangeWorkStateLogic notify json Marshal uid:%d err:%+v", in.ListenerUid, err)
		return
	}

	err = l.svcCtx.KqSubscribeNotifyMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("ChangeWorkStateLogic notify kafka Push uid:%d err:%+v", in.ListenerUid, err)
		return
	}
}
