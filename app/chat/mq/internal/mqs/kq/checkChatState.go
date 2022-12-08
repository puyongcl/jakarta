package kq

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/chat/mq/internal/svc"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/kqueue"
	"time"
)

type CheckChatStateMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckChatStateMq(ctx context.Context, svcCtx *svc.ServiceContext) *CheckChatStateMq {
	return &CheckChatStateMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckChatStateMq) Consume(_, val string) error {
	var message kqueue.CheckChatStateMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("CheckChatStateMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("CheckChatStateMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *CheckChatStateMq) execService(message *kqueue.CheckChatStateMessage) error {
	msg := jobtype.DeferCheckChatStatePayload{
		Uid:         message.Uid,
		ListenerUid: message.ListenerUid,
		OrderType:   message.OrderType,
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferCheckChatState, buf), asynq.ProcessIn(time.Duration(message.DeferMinute)*time.Minute))
	if err != nil {
		return err
	}
	return nil
}
