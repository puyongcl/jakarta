package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/statPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/userkey"
	"time"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatisticDailyDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStatisticDailyDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatisticDailyDataLogic {
	return &UpdateStatisticDailyDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新每日统计数据
func (l *UpdateStatisticDailyDataLogic) UpdateStatisticDailyData(in *pb.UpdateStatisticDailyDataReq) (*pb.UpdateStatisticDailyDataResp, error) {
	// 获取时间范围
	var start, end time.Time
	var err error
	start, err = time.ParseInLocation(db.DateTimeFormat, in.StartTime, time.Local)
	if err != nil {
		return nil, err
	}
	end, err = time.ParseInLocation(db.DateTimeFormat, in.EndTime, time.Local)
	if err != nil {
		return nil, err
	}
	// 查询渠道列表
	var chanList []string

	chanList, err = l.svcCtx.UserAuthModel.FindChannelList(l.ctx, &start, &end)
	if err != nil {
		return nil, err
	}
	chanList = append(chanList, "")
	// 查询数据
	dateStr := time.Now().AddDate(0, 0, -1).Format(db.DateFormat2)

	for k, _ := range chanList {
		var data *statPgModel.StatDaily
		data, err = l.queryStatData(&start, &end, chanList[k], dateStr)
		if err != nil {
			return nil, err
		}

		_, err = l.svcCtx.StatDailyModel.Insert(l.ctx, data)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UpdateStatisticDailyDataLogic StatDailyModel.Insert channel:%+v err:%+v", chanList, err)
			err = nil
		}
	}
	return &pb.UpdateStatisticDailyDataResp{}, nil
}

func (l *UpdateStatisticDailyDataLogic) queryStatData(start, end *time.Time, channel, dateStr string) (data *statPgModel.StatDaily, err error) {
	data = new(statPgModel.StatDaily)
	// 统计昨天
	if channel != "" {
		data.Id = fmt.Sprintf(db.DBUidId2, dateStr, channel)
	} else {
		data.Id = dateStr
	}
	data.Channel = channel
	data.CreateTime = time.Now()
	data.NewUserCnt, err = l.svcCtx.UserAuthModel.CountNewUser(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.LoginUserCnt, err = l.svcCtx.UserLoginLogModel.Count(l.ctx, start, end, channel, userkey.UserTypeNormalUser)
	if err != nil {
		return nil, err
	}
	data.LoginListenerCnt, err = l.svcCtx.UserLoginLogModel.Count(l.ctx, start, end, channel, userkey.UserTypeListener)
	if err != nil {
		return nil, err
	}
	data.PaidUserCnt, err = l.svcCtx.ChatOrderModel.CountPaidUserCnt(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.PaidOrderCnt, err = l.svcCtx.ChatOrderModel.CountPaidOrderCnt(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.PaidAmount, err = l.svcCtx.ChatOrderModel.SumPaidAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.ApplyRefundAmount, err = l.svcCtx.ChatOrderModel.SumApplyRefundAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.RefundSuccessAmount, err = l.svcCtx.ChatOrderModel.SumRefundSuccessAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.ConfirmOrderAmount, err = l.svcCtx.ChatOrderModel.SumConfirmAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.ListenerAmount, err = l.svcCtx.ChatOrderModel.SumListenerAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	data.PlatformAmount, err = l.svcCtx.ChatOrderModel.SumPlatformAmount(l.ctx, start, end, channel)
	if err != nil {
		return nil, err
	}
	return
}
