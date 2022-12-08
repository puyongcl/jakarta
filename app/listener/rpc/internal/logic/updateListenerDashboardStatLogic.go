package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/dashboard"
	"jakarta/common/key/db"
	"jakarta/common/xerr"
	"strconv"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerDashboardStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerDashboardStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerDashboardStatLogic {
	return &UpdateListenerDashboardStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX首页数据统计看板
func (l *UpdateListenerDashboardStatLogic) UpdateListenerDashboardStat(in *pb.UpdateListenerDashboardStatReq) (*pb.UpdateListenerDashboardStatResp, error) {
	if len(in.ListenerUid) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty listener uid")
	}
	var err error
	var now time.Time
	now, err = time.ParseInLocation(db.DateTimeFormat, in.Time, time.Local)
	if err != nil {
		return nil, err
	}
	var rank, score int64

	mm := make(map[string]string, 0)
	var rs *listenerPgModel2.ListenerProfile
	var data *dashboard.ListenerDashboardRedisHashData
	var lastUpdateTime time.Time
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		//
		data, err = l.svcCtx.StatRedis.GetListenerDashboard(l.ctx, in.ListenerUid[idx])
		if err != nil {
			return nil, err
		}
		// 如果是今天第一次更新 则记录昨天的数据
		lastUpdateTime, err = time.ParseInLocation(db.DateTimeFormat, data.TodayStatUpdateTime, time.Local)
		if err != nil {
			return nil, err
		}

		if lastUpdateTime.Day() != now.Day() { // 今天第一次更新
			err = l.snapshot(data)
			if err != nil {
				return nil, err
			}
		}

		// 今日订单数和排名
		rank, score, err = l.svcCtx.StatRedis.GetListenerTodayOrderCntAndRank(l.ctx, in.ListenerUid[idx])
		if err != nil && err != redis.Nil {
			return nil, err
		}
		mm["TodayOrderCnt"] = strconv.FormatInt(score, 10)
		mm["TodayOrderCntRank"] = strconv.FormatInt(rank, 10)

		// 今日订单金额
		rank, score, err = l.svcCtx.StatRedis.GetListenerTodayOrderAmountAndRank(l.ctx, in.ListenerUid[idx])
		if err != nil && err != redis.Nil {
			return nil, err
		}
		mm["TodayOrderAmount"] = strconv.FormatInt(score, 10)
		mm["TodayOrderAmountRank"] = strconv.FormatInt(rank, 10)

		// 今日推荐人数和排名
		rank, score, err = l.svcCtx.StatRedis.GetListenerTodayRecommendUserCntAndRank(l.ctx, in.ListenerUid[idx])
		if err != nil && err != redis.Nil {
			return nil, err
		}
		mm["TodayRecommendUserCnt"] = strconv.FormatInt(score, 10)
		mm["TodayRecommendUserCntRank"] = strconv.FormatInt(rank, 10)

		// 今日进入聊天页面人数和排名
		rank, score, err = l.svcCtx.StatRedis.GetListenerTodayEnterChatCntAndRank(l.ctx, in.ListenerUid[idx])
		if err != nil && err != redis.Nil {
			return nil, err
		}
		mm["TodayEnterChatUserCnt"] = strconv.FormatInt(score, 10)
		mm["TodayEnterChatUserCntRank"] = strconv.FormatInt(rank, 10)

		// 今日访问个人资料页面人数和排名
		rank, score, err = l.svcCtx.StatRedis.GetListenerTodayViewUserCntAndRank(l.ctx, in.ListenerUid[idx])
		if err != nil && err != redis.Nil {
			return nil, err
		}
		mm["TodayViewUserCnt"] = strconv.FormatInt(score, 10)
		mm["TodayViewUserCntRank"] = strconv.FormatInt(rank, 10)

		// 累计差评和退款数
		rs, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid[idx])
		if err != nil {
			return nil, err
		}
		mm["OneStarRatingOrderCnt"] = strconv.FormatInt(rs.OneStar, 10)
		mm["RefundOrderCnt"] = strconv.FormatInt(rs.RefundOrderCnt, 10)

		// 更新时间
		mm["TodayStatUpdateTime"] = in.Time

		err = l.svcCtx.StatRedis.HMSetListenerDashboard(l.ctx, in.ListenerUid[idx], mm)
		if err != nil {
			return nil, err
		}
	}

	return &pb.UpdateListenerDashboardStatResp{}, nil
}

// 保存昨日统计的每日数据
func (l *UpdateListenerDashboardStatLogic) snapshot(in *dashboard.ListenerDashboardRedisHashData) error {
	var data dashboard.ListenerDashboard
	err := dashboard.TransferListenerDashboardData(in, &data)
	if err != nil {
		return err
	}
	var stat *listenerPgModel2.ListenerDashboardStat
	stat, err = l.svcCtx.ListenerDashboardStatModel.FindOne(l.ctx, data.ListenerUid)
	if err != nil {
		return err
	}
	_ = copier.Copy(stat, &data)
	stat.YesterdayOrderCnt = data.TodayOrderCnt
	stat.YesterdayOrderCntRank = data.TodayOrderCntRank
	stat.YesterdayOrderAmount = data.TodayOrderAmount
	stat.YesterdayOrderAmountRank = data.TodayOrderAmountRank
	stat.YesterdayRecommendUserCnt = data.TodayRecommendUserCnt
	stat.YesterdayRecommendUserCntRank = data.TodayRecommendUserCntRank
	stat.YesterdayEnterChatUserCnt = data.TodayEnterChatUserCnt
	stat.YesterdayEnterChatUserCntRank = data.TodayEnterChatUserCntRank
	stat.YesterdayViewUserCnt = data.TodayViewUserCnt
	stat.YesterdayViewUserCntRank = data.TodayViewUserCntRank

	stat.YesterdayLastUpdateTime, err = time.ParseInLocation(db.DateTimeFormat, data.TodayStatUpdateTime, time.Local)
	if err != nil {
		return err
	}
	stat.LastDayLastUpdateTime, err = time.ParseInLocation(db.DateTimeFormat, data.LastDayStatUpdateTime, time.Local)
	if err != nil {
		return err
	}
	err = l.svcCtx.ListenerDashboardStatModel.Update(l.ctx, stat)
	if err != nil {
		return err
	}
	return nil
}
