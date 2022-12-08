package listener

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"time"
)

type AutoUpdateListenerDashboardHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoUpdateListenerDashboardHandler(svcCtx *svc.ServiceContext) *AutoUpdateListenerDashboardHandler {
	return &AutoUpdateListenerDashboardHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoUpdateListenerDashboardHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取最近活跃的XXX
	var pageNo int64 = 1
	now := time.Now()
	end := now.Format(db.DateTimeFormat)
	start := now.AddDate(0, 0, -listenerkey.AutoUpdateListenerUserStatRangeDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbListener.FindListenerListRangeByUpdateTimeResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.ListenerRpc.FindListenerListRangeByUpdateTime(ctx, &pbListener.FindListenerListRangeByUpdateTimeReq{
			PageNo:   pageNo,
			PageSize: 10,
			Start:    start,
			End:      end,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoUpdateListenerDashboardHandler FindListenerListRangeByUpdateTime err:%+v", err)
			return err
		}

		if len(rsp.Listener) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoUpdateListenerDashboardHandler exit. process listener cnt:%d", cnt)
			}
			break
		}

		// 更新XXX首页数据看板
		var in1 pbListener.UpdateListenerDashboardStatReq
		in1.ListenerUid = make([]int64, 0)
		in1.Time = now.Format(db.DateTimeFormat)
		for idx := 0; idx < len(rsp.Listener); idx++ {
			in1.ListenerUid = append(in1.ListenerUid, rsp.Listener[idx].ListenerUid)
		}
		_, err = l.svcCtx.ListenerRpc.UpdateListenerDashboardStat(ctx, &in1)
		if err != nil {
			return err
		}

		cnt += len(rsp.Listener)
	}
	logx.WithContext(ctx).Infof("AutoUpdateListenerDashboardHandler done")
	return nil
}
