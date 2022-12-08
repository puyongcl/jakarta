package stat

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/statistic/rpc/pb"
)

type AutoUpdateUserStatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoUpdateUserStatHandler(svcCtx *svc.ServiceContext) *AutoUpdateUserStatHandler {
	return &AutoUpdateUserStatHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoUpdateUserStatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	_, err := l.svcCtx.StatRpc.UpdateUserStateStat(ctx, &pb.UpdateUserStateStatReq{})
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("AutoUpdateUserStatHandler done")
	return nil
}
