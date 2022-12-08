package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/statPgModel"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"time"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDailyStatListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDailyStatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDailyStatListLogic {
	return &GetDailyStatListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取每日统计数据
func (l *GetDailyStatListLogic) GetDailyStatList(in *pb.GetDailyStatListReq) (*pb.GetDailyStatListResp, error) {
	resp := pb.GetDailyStatListResp{List: make([]*pb.DailyStat, 0)}
	var id string
	todayId := time.Now().Format(db.DateFormat2)
	var rsp *statPgModel.StatDaily
	var err error
	for idx := 0; idx < len(in.Date); idx++ {
		id = in.Date[idx]
		if in.Channel != "" {
			id = fmt.Sprintf(db.DBUidId2, in.Date[idx], in.Channel)
		}
		if in.Date[idx] == todayId { // 今日数据 实时查询
			usd := NewUpdateStatisticDailyDataLogic(l.ctx, l.svcCtx)
			ts, te := tool.GetTodayStartAndEndTime()
			rsp, err = usd.queryStatData(ts, te, in.Channel, todayId)
			if err != nil {
				return nil, err
			}
		} else {
			rsp, err = l.svcCtx.StatDailyModel.FindOne(l.ctx, id)
			if err != nil && err != statPgModel.ErrNotFound {
				return nil, err
			}
		}

		if rsp != nil {
			var val pb.DailyStat
			_ = copier.Copy(&val, rsp)
			val.CreateTime = rsp.CreateTime.Format(db.DateTimeFormat)
			resp.List = append(resp.List, &val)
		}
	}
	return &resp, nil
}
