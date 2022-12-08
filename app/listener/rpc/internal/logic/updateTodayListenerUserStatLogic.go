package logic

import (
	"context"
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTodayListenerUserStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTodayListenerUserStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodayListenerUserStatLogic {
	return &UpdateTodayListenerUserStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新统计今日推荐和浏览XXX资料页统计
func (l *UpdateTodayListenerUserStatLogic) UpdateTodayListenerUserStat(in *pb.UpdateTodayListenerUserStatReq) (*pb.UpdateTodayListenerUserStatResp, error) {
	if len(in.ListenerUid) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty listener uid")
	}

	start, end := tool.GetTodayStartAndEndTime()

	var cnt int64
	var err error
	// 更新今日推荐用户数
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		cnt, err = l.svcCtx.ListenerUserRecommendStatModel.FindCountRecommendUserCntRangeCreateTime(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.StatRedis.SetTodayListenerRecommendUserCnt(l.ctx, in.ListenerUid[idx], cnt)
		if err != nil {
			return nil, err
		}
	}

	// 更新今日浏览用户数
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		cnt, err = l.svcCtx.ListenerUserViewStatModel.FindCountViewUserCntRangeCreateTime(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.StatRedis.SetTodayListenerViewUserCnt(l.ctx, in.ListenerUid[idx], cnt)
		if err != nil {
			return nil, err
		}
	}
	return &pb.UpdateTodayListenerUserStatResp{}, nil
}
