package logic

import (
	"context"
	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"
	"jakarta/common/dashboard"
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLastDaysEnterChatUserCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLastDaysEnterChatUserCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLastDaysEnterChatUserCntLogic {
	return &UpdateLastDaysEnterChatUserCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新统计近几天进入XXX页面用户数
func (l *UpdateLastDaysEnterChatUserCntLogic) UpdateLastDaysEnterChatUserCnt(in *pb.UpdateLastDaysEnterChatUserCntReq) (*pb.UpdateLastDaysEnterChatUserCntResp, error) {
	if len(in.ListenerUid) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty listener uid")
	}

	start, end := tool.GetLastDayStartAndEndTime(dashboard.LastStat30Days)

	var cnt int64
	var err error

	// 更新最近30天进入XXX聊天界面的用户数
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		cnt, err = l.svcCtx.ChatBalanceModel.FindCountEnterChatUserCntRangeUpdateTime(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ChatRedis.SetListenerDashboard(l.ctx, in.ListenerUid[idx], "Last30DaysEnterChatUserCnt", cnt)
		if err != nil {
			return nil, err
		}
	}

	start2 := end.AddDate(0, 0, -dashboard.LastStat7Days)
	// 更新最近7天进入XXX聊天界面的用户数
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		cnt, err = l.svcCtx.ChatBalanceModel.FindCountEnterChatUserCntRangeUpdateTime(l.ctx, in.ListenerUid[idx], &start2, end)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ChatRedis.SetListenerDashboard(l.ctx, in.ListenerUid[idx], "Last7DaysEnterChatUserCnt", cnt)
		if err != nil {
			return nil, err
		}
	}
	return &pb.UpdateLastDaysEnterChatUserCntResp{}, nil
}
