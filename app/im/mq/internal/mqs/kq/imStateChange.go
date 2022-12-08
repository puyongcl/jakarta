package kq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	pbListener "jakarta/app/listener/rpc/pb"
	pbStat "jakarta/app/statistic/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/kqueue"
)

// IM回调消息 用户登陆下线状态处理
type ImStateChangeMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImCallbackMq(ctx context.Context, svcCtx *svc.ServiceContext) *ImStateChangeMq {
	return &ImStateChangeMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImStateChangeMq) Consume(_, val string) error {
	var message kqueue.ImStateChangeMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("ImStateChangeMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("ImStateChangeMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *ImStateChangeMq) execService(message *kqueue.ImStateChangeMessage) error {
	rsp, err := l.svcCtx.UserRpc.UpdateUserLoginState(l.ctx, &pbUser.UpdateUserLoginStateReq{
		Uid:       message.Uid,
		State:     message.State,
		EventTime: message.EventTime,
	})
	if err != nil {
		return err
	}

	if rsp.UserType == userkey.UserTypeListener && rsp.IsUpdated > 0 {
		_, err = l.svcCtx.ListenerRpc.UpdateListenerOnlineState(l.ctx, &pbListener.UpdateListenerOnlineStateReq{
			State:         message.State,
			ListenerUid:   message.Uid,
			TodayLoginCnt: rsp.TodayLoginCnt,
		})
		if err != nil {
			return err
		}
	}
	_, err = l.svcCtx.StatRpc.UpdateLoginLog(l.ctx, &pbStat.UpdateLoginLogReq{Uid: message.Uid, State: message.State, UserType: rsp.UserType, Channel: rsp.Channel})
	if err != nil {
		return err
	}
	return nil
}
