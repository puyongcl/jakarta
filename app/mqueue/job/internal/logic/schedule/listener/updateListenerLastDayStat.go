package listener

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"time"
)

type UpdateListenerLastDayStatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateListenerLastDayStatHandler(svcCtx *svc.ServiceContext) *UpdateListenerLastDayStatHandler {
	return &UpdateListenerLastDayStatHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *UpdateListenerLastDayStatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 保存XXX最近几日统计数据
	err := l.updateListenerLastDayStat(ctx)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	// 更新XXX的统计平均数据
	var in4 pbListener.UpdateListenerEveryDayAverageStatReq
	_, err = l.svcCtx.ListenerRpc.UpdateListenerEveryDayAverageStat(ctx, &in4)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	// 更新XXX建议
	err = l.updateListenerSuggestion(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l *UpdateListenerLastDayStatHandler) updateListenerLastDayStat(ctx context.Context) error {
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
			logx.WithContext(ctx).Errorf("UpdateListenerLastDayStatHandler updateListenerLastDayStat FindListenerListRangeByUpdateTime err:%+v", err)
			return err
		}

		if len(rsp.Listener) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("UpdateListenerLastDayStatHandler updateListenerLastDayStat exit. process order cnt:%d", cnt)
			}
			break
		}
		var in1 pbChat.UpdateLastDaysEnterChatUserCntReq
		var in2 pbOrder.UpdateOrderLastDaysStatReq
		var in3 pbListener.SnapshotLastDaysListenerStatReq
		in1.ListenerUid = make([]int64, 0)
		in2.ListenerUid = make([]int64, 0)
		in3.ListenerUid = make([]int64, 0)

		for idx := 0; idx < len(rsp.Listener); idx++ {
			in1.ListenerUid = append(in1.ListenerUid, rsp.Listener[idx].ListenerUid)
			in2.ListenerUid = append(in2.ListenerUid, rsp.Listener[idx].ListenerUid)
			in3.ListenerUid = append(in3.ListenerUid, rsp.Listener[idx].ListenerUid)
		}

		// 更新聊天进入用户数
		_, err = l.svcCtx.ChatRpc.UpdateLastDaysEnterChatUserCnt(ctx, &in1)
		if err != nil {
			return err
		}
		// 更新最近几天订单统计数据
		_, err = l.svcCtx.OrderRpc.UpdateOrderLastDaysStat(ctx, &in2)
		if err != nil {
			return err
		}

		// 保存昨天的统计数据
		_, err = l.svcCtx.ListenerRpc.SnapshotLastDaysListenerStat(ctx, &in3)
		if err != nil {
			return err
		}

		cnt += len(rsp.Listener)
	}
	logx.WithContext(ctx).Infof("UpdateListenerLastDayStatHandler done")
	return nil
}

func (l *UpdateListenerLastDayStatHandler) updateListenerSuggestion(ctx context.Context) error {
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
			logx.WithContext(ctx).Errorf("UpdateListenerLastDayStatHandler updateListenerSuggestion FindListenerListRangeByUpdateTime err:%+v", err)
			return err
		}

		if len(rsp.Listener) <= 0 {
			logx.WithContext(ctx).Infof("UpdateListenerLastDayStatHandler updateListenerSuggestion exit. process listener cnt:%d", cnt)
			break
		}

		// 更新XXX的每日建议
		var in1 pbListener.UpdateListenerSuggestionReq
		in1.ListenerUid = make([]int64, 0)
		for idx := 0; idx < len(rsp.Listener); idx++ {
			in1.ListenerUid = append(in1.ListenerUid, rsp.Listener[idx].ListenerUid)
		}
		_, err = l.svcCtx.ListenerRpc.UpdateListenerSuggestion(ctx, &in1)
		if err != nil {
			return err
		}

		cnt += len(rsp.Listener)
	}
	return nil
}
