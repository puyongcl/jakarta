package listener

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/common/key/rediskey"
)

type AutoClearHistoryDataHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoClearHistoryDataHandler(svcCtx *svc.ServiceContext) *AutoClearHistoryDataHandler {
	return &AutoClearHistoryDataHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoClearHistoryDataHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 清理XXX每日排行历史数据
	err := l.svcCtx.JobRedis.ClearTodayStatHistoryData(ctx, rediskey.NeedDelTodayRedisData)
	if err != nil {
		return err
	}

	// 清理新用户推荐XXX数据
	err = l.svcCtx.JobRedis.ClearRecommendListenerPoolData(ctx, rediskey.NeedDelRedisRecommendData)
	if err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("AutoClearHistoryDataHandler done")
	return nil
}
