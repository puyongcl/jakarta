package logic

import (
	"context"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/average"
	"jakarta/common/key/listenerkey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerSuggestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerSuggestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerSuggestionLogic {
	return &UpdateListenerSuggestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新XXX建议
func (l *UpdateListenerSuggestionLogic) UpdateListenerSuggestion(in *pb.UpdateListenerSuggestionReq) (*pb.UpdateListenerSuggestionResp, error) {
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		err := l.updateSuggestion(in.ListenerUid[idx])
		if err != nil {
			return nil, err
		}
	}

	return &pb.UpdateListenerSuggestionResp{}, nil
}

func (l *UpdateListenerSuggestionLogic) updateSuggestion(listenerUid int64) error {
	var err error
	var lp *listenerPgModel2.ListenerProfile
	lp, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, listenerUid)
	if err != nil {
		return err
	}
	var sg []int64
	if lp.PaidOrderCnt <= 10 {
		sg = append(sg, listenerkey.SuggestionNo1)
		return nil
	} else {
		var avg *average.ListenerStatAverage
		avg, err = l.svcCtx.StatRedis.GetListenerAverageStat(l.ctx)
		if err != nil {
			return err
		}
		// 昨日数据
		var lastStat *listenerPgModel2.ListenerDashboardStat
		lastStat, err = l.svcCtx.ListenerDashboardStatModel.FindOne(l.ctx, listenerUid)
		if err != nil {
			return err
		}
		// 昨日数据 与 昨日平均比较
		// 昨日曝光量、点击量低于平均值
		if lastStat.YesterdayRecommendUserCnt < avg.YesterdayRecommendUserCnt && lastStat.YesterdayEnterChatUserCnt < avg.YesterdayEnterChatUserCnt {
			sg = append(sg, listenerkey.SuggestionNo2)
		}
		// 昨日近多少天数据 与 昨天近多少天平均比较
		// 近30天下单率低于平均值
		if lastStat.Last30DaysPaidUserRate < avg.Last30DaysPaidUserRate {
			sg = append(sg, listenerkey.SuggestionNo3)
		}
		// 近30天复购率低于平均值
		if lastStat.Last30DaysRepeatPaidUserRate < avg.Last30DaysRepeatPaidUserRate {
			sg = append(sg, listenerkey.SuggestionNo4)
		}
		// 近7天订单金额低于平均值
		if lastStat.Last7DaysPaidAmountSum < avg.Last7DaysPaidAmountSum {
			sg = append(sg, listenerkey.SuggestionNo5)
		}
		// 如果全部不符合
		if len(sg) <= 0 {
			sg = append(sg, listenerkey.SuggestionNo6)
		}
	}

	err = l.svcCtx.ListenerDashboardStatModel.UpdateSuggestion(l.ctx, listenerUid, sg)
	if err != nil {
		return err
	}
	return nil
}
