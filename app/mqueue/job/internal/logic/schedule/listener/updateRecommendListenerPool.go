package listener

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
)

type UpdateRecommendListenerPoolHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateNewUserRecommendListenerHandler(svcCtx *svc.ServiceContext) *UpdateRecommendListenerPoolHandler {
	return &UpdateRecommendListenerPoolHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *UpdateRecommendListenerPoolHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 更新今日推荐XXX列表
	in := pbListener.UpdateRecommendListenerPoolReq{
		RecentDay: 3,
		Size:      300,
	}

	if l.svcCtx.Config.Mode == service.DevMode {
		in.RecentDay = 100
	}

	rsp, err := l.svcCtx.ListenerRpc.UpdateRecommendListenerPool(ctx, &in)
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("UpdateRecommendListenerPoolHandler done cnt:%d", rsp.Cnt)
	return nil
}
