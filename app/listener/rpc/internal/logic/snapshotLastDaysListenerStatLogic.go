package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/dashboard"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SnapshotLastDaysListenerStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSnapshotLastDaysListenerStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SnapshotLastDaysListenerStatLogic {
	return &SnapshotLastDaysListenerStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  保存最近多少天的统计数据（一天更新一次，不能覆盖每日更新的数据)
func (l *SnapshotLastDaysListenerStatLogic) SnapshotLastDaysListenerStat(in *pb.SnapshotLastDaysListenerStatReq) (*pb.SnapshotLastDaysListenerStatResp, error) {
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		err := l.snapshot(in.ListenerUid[idx])
		if err != nil {
			return nil, err
		}
	}

	return &pb.SnapshotLastDaysListenerStatResp{}, nil
}

func (l *SnapshotLastDaysListenerStatLogic) snapshot(listenerUid int64) error {
	statData, err := l.svcCtx.StatRedis.GetListenerDashboard(l.ctx, listenerUid)
	if err != nil {
		return err
	}
	var data dashboard.ListenerDashboard
	err = dashboard.TransferListenerDashboardData(statData, &data)
	if err != nil {
		return err
	}
	var stat listenerPgModel.ListenerDashboardStatLastDays
	_ = copier.Copy(&stat, &data)
	stat.LastDayLastUpdateTime, err = time.ParseInLocation(db.DateTimeFormat, data.LastDayStatUpdateTime, time.Local)
	if err != nil {
		return err
	}
	// 计算比例
	stat.Last30DaysPaidUserRate = tool.DivideInt64(stat.Last30DaysPaidUserCnt*dashboard.DivideNumber, stat.Last30DaysEnterChatUserCnt)
	stat.Last30DaysRepeatPaidUserRate = tool.DivideInt64(stat.Last30DaysRepeatPaidUserCnt*dashboard.DivideNumber, stat.Last30DaysPaidUserCnt)
	stat.Last7DaysPaidUserRate = tool.DivideInt64(stat.Last7DaysPaidUserCnt*dashboard.DivideNumber, stat.Last7DaysEnterChatUserCnt)
	stat.Last7DaysRepeatPaidUserRate = tool.DivideInt64(stat.Last7DaysRepeatPaidUserCnt*dashboard.DivideNumber, stat.Last30DaysPaidUserCnt)
	err = l.svcCtx.ListenerDashboardStatModel.UpdateLastDaysStat(l.ctx, listenerUid, &stat)
	if err != nil {
		return err
	}
	return nil
}
