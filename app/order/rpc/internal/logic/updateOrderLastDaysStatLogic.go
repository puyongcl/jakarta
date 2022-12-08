package logic

import (
	"context"
	"jakarta/common/dashboard"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strconv"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLastDaysStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderLastDaysStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLastDaysStatLogic {
	return &UpdateOrderLastDaysStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX订单统计数据
func (l *UpdateOrderLastDaysStatLogic) UpdateOrderLastDaysStat(in *pb.UpdateOrderLastDaysStatReq) (*pb.UpdateOrderLastDaysStatResp, error) {
	if len(in.ListenerUid) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty listener uid")
	}

	start, end := tool.GetLastDayStartAndEndTime(dashboard.LastStat30Days)
	var amount int64
	var cnt int64
	var err error
	var val int64

	// 更新最近30天
	mm := make(map[string]string, 0)
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		// 最近30天下单人数
		cnt, err = l.svcCtx.ChatOrderModel.CountPaidUser(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm["Last30DaysPaidUserCnt"] = strconv.FormatInt(cnt, 10)

		// 过去30天总消费
		amount, err = l.svcCtx.ChatOrderModel.SumListenerPaidAmount(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm["Last30DaysPaidAmountSum"] = strconv.FormatInt(amount, 10)

		// 过去30天人均消费
		val = tool.DivideInt64(amount, cnt)
		mm["Last30DaysAveragePaidAmountPerUser"] = strconv.FormatInt(val, 10)

		// 过去30天日均消费
		val = tool.DivideInt64(amount, dashboard.LastStat30Days)
		mm["Last30DaysAveragePaidAmountPerDay"] = strconv.FormatInt(val, 10)

		// 过去30天复购人数
		cnt, err = l.svcCtx.ChatOrderModel.CountRepeatPaidUser(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm["Last30DaysRepeatPaidUserCnt"] = strconv.FormatInt(cnt, 10)
		err = l.svcCtx.StatRedis.HMSetListenerDashboard(l.ctx, in.ListenerUid[idx], mm)
		if err != nil {
			return nil, err
		}
	}

	// 过去7天
	t := end.AddDate(0, 0, -dashboard.LastStat7Days)
	start = &t
	now := time.Now().Format(db.DateTimeFormat)
	mm2 := make(map[string]string, 0)
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		// 最近7天下单人数
		cnt, err = l.svcCtx.ChatOrderModel.CountPaidUser(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm2["Last7DaysPaidUserCnt"] = strconv.FormatInt(cnt, 10)

		// 过去7天总消费
		amount, err = l.svcCtx.ChatOrderModel.SumListenerPaidAmount(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm["Last7DaysPaidAmountSum"] = strconv.FormatInt(amount, 10)

		// 过去7天人均消费
		val = tool.DivideInt64(amount, cnt)
		mm2["Last7DaysAveragePaidAmountPerUser"] = strconv.FormatInt(val, 10)

		// 过去7天日均消费
		val = tool.DivideInt64(amount, dashboard.LastStat7Days)
		mm2["Last7DaysAveragePaidAmountPerDay"] = strconv.FormatInt(val, 10)
		// 过去7天复购人数
		cnt, err = l.svcCtx.ChatOrderModel.CountRepeatPaidUser(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		mm2["Last7DaysRepeatPaidUserCnt"] = strconv.FormatInt(cnt, 10)

		// 更新时间
		mm2["LastDayStatUpdateTime"] = now
		err = l.svcCtx.StatRedis.HMSetListenerDashboard(l.ctx, in.ListenerUid[idx], mm2)
		if err != nil {
			return nil, err
		}
	}

	return &pb.UpdateOrderLastDaysStatResp{}, nil
}
