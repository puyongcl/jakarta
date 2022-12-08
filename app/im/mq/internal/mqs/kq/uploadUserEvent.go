package kq

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	"jakarta/common/kqueue"
)

// 上传渠道用户事件
type UploadUserEventMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadUserEventMq(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUserEventMq {
	return &UploadUserEventMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadUserEventMq) Consume(_, val string) error {
	var msg kqueue.UploadUserEventMessage
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		logx.WithContext(l.ctx).Errorf("UploadUserEventMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&msg); err != nil {
		logx.WithContext(l.ctx).Errorf("UploadUserEventMq->execService err : %+v , val : %s , msg:%+v", err, val, msg)
		return err
	}

	return nil
}

func (l *UploadUserEventMq) execService(message *kqueue.UploadUserEventMessage) error {
	return l.svcCtx.Zrc.UploadEvent(message.Cb, message.Event, message.Value, message.Stamp)
}
