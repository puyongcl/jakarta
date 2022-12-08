package kq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/payment/mq/internal/svc"
	"jakarta/common/kqueue"
)

type CommitMoveCashMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommitMoveCashMq(ctx context.Context, svcCtx *svc.ServiceContext) *CommitMoveCashMq {
	return &CommitMoveCashMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommitMoveCashMq) Consume(_, val string) error {
	var message kqueue.CommitMoveCashMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("CommitMoveCashMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("CommitMoveCashMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}
	return nil
}

// TODO 废弃
func (l *CommitMoveCashMq) execService(message *kqueue.CommitMoveCashMessage) error {
	logx.WithContext(l.ctx).Infof("CommitMoveCashMq message:%+v", message)

	return nil
}
