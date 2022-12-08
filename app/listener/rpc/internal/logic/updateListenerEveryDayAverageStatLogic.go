package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/average"
	"jakarta/common/key/db"
	"reflect"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerEveryDayAverageStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerEveryDayAverageStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerEveryDayAverageStatLogic {
	return &UpdateListenerEveryDayAverageStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX每日统计数据的平均值
func (l *UpdateListenerEveryDayAverageStatLogic) UpdateListenerEveryDayAverageStat(in *pb.UpdateListenerEveryDayAverageStatReq) (*pb.UpdateListenerEveryDayAverageStatResp, error) {
	now := time.Now()
	var err error
	id := now.Format(db.DateFormat2)
	var rs *listenerPgModel2.ListenerStatAverage
	rs, err = l.svcCtx.ListenerStatAverageModel.FindOne(l.ctx, id)
	if err != nil && err != listenerPgModel2.ErrNotFound {
		return nil, err
	}
	err = nil
	if rs != nil {
		logx.WithContext(l.ctx).Errorf("UpdateListenerEveryDayAverageStatLogic 已经执行过定时任务")
		return &pb.UpdateListenerEveryDayAverageStatResp{}, nil
	}
	// 获取平均值
	var data listenerPgModel2.ListenerStatAverage
	t1 := reflect.TypeOf(listenerPgModel2.ListenerDashboardStat{})
	t2 := reflect.TypeOf(data)
	tv2 := reflect.ValueOf(&data)
	mv := make(map[string]int64, 0)
	var cnt float64
	var dbFieldName, fieldName string
	for k := 0; k < t1.NumField(); k++ {
		if t1.Field(k).Type.Kind() == reflect.Int64 {
			fieldName = t1.Field(k).Name
			if fieldName == "ListenerUid" {
				continue
			}
			dbFieldName = t1.Field(k).Tag.Get("db")
			cnt, err = l.svcCtx.ListenerDashboardStatModel.CalculateAverage(l.ctx, dbFieldName)
			if err != nil {
				return nil, err
			}
			mv[fieldName] = int64(cnt)
		}
	}

	for k := 0; k < t2.NumField(); k++ {
		if t2.Field(k).Type.Kind() == reflect.Int64 {
			fieldName = t2.Field(k).Name
			v, ok := mv[fieldName]
			if ok {
				tv2.Elem().FieldByName(fieldName).SetInt(v)
			}
		}
	}

	//
	data.Id = id
	_, err = l.svcCtx.ListenerStatAverageModel.Insert(l.ctx, &data)
	if err != nil {
		return nil, err
	}
	//
	var cData average.ListenerStatAverage
	_ = copier.Copy(&cData, data)
	cData.CreateTime = now.Format(db.DateTimeFormat)
	err = l.svcCtx.StatRedis.SetListenerAverageStat(l.ctx, &cData)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateListenerEveryDayAverageStatResp{}, nil
}
