package delay

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/mqueue/job/jobtype"
)

// DeferSendImMsgHandler close no pay chatOrder
type DeferSendImMsgHandler struct {
	svcCtx *svc.ServiceContext
}

func NewDeferSendImMsgHandler(svcCtx *svc.ServiceContext) *DeferSendImMsgHandler {
	return &DeferSendImMsgHandler{
		svcCtx: svcCtx,
	}
}

// defer send im msg
func (l *DeferSendImMsgHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p jobtype.DeferSendImMsgPayload
	err := json.Unmarshal(t.Payload(), &p)
	if err != nil {
		return err
	}

	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(p.KqMsgBuf))
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("DeferSendImMsgHandler done")
	return nil
}
