package stat

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/statistic/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/tool"
)

type UpdateDailyStatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateDailyStatHandler(svcCtx *svc.ServiceContext) *UpdateDailyStatHandler {
	return &UpdateDailyStatHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *UpdateDailyStatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	start, end := tool.GetYesterdayStartAndEndTime()
	in := pb.UpdateStatisticDailyDataReq{
		StartTime: start.Format(db.DateTimeFormat),
		EndTime:   end.Format(db.DateTimeFormat),
	}
	_, err := l.svcCtx.StatRpc.UpdateStatisticDailyData(ctx, &in)
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("AutoUpdateUserStatHandler done")
	return nil
}
