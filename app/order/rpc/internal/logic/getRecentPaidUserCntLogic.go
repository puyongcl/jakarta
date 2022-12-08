package logic

import (
	"context"
	"jakarta/common/key/db"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentPaidUserCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecentPaidUserCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecentPaidUserCntLogic {
	return &GetRecentPaidUserCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取最近时间段付费用户数
func (l *GetRecentPaidUserCntLogic) GetRecentPaidUserCnt(in *pb.GetRecentPaidUserCntReq) (*pb.GetRecentPaidUserCntResp, error) {
	var start, end time.Time
	start, err := time.ParseInLocation(db.DateTimeFormat, in.StartTime, time.Local)
	if err != nil {
		return nil, err
	}
	end, err = time.ParseInLocation(db.DateTimeFormat, in.EndTime, time.Local)
	if err != nil {
		return nil, err
	}
	resp := pb.GetRecentPaidUserCntResp{}
	resp.UserCnt, err = l.svcCtx.ChatOrderModel.CountPaidUserCnt2(l.ctx, in.ListenerUid, &start, &end)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
