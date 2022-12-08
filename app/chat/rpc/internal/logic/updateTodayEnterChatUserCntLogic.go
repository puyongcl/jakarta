package logic

import (
	"context"
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTodayEnterChatUserCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTodayEnterChatUserCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodayEnterChatUserCntLogic {
	return &UpdateTodayEnterChatUserCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新统计进入XXX聊天页面用户数
func (l *UpdateTodayEnterChatUserCntLogic) UpdateTodayEnterChatUserCnt(in *pb.UpdateTodayEnterChatUserCntReq) (*pb.UpdateTodayEnterChatUserCntResp, error) {
	if len(in.ListenerUid) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty listener uid")
	}

	start, end := tool.GetTodayStartAndEndTime()

	var cnt int64
	var err error
	// 更新今日进入XXX聊天界面的用户数
	for idx := 0; idx < len(in.ListenerUid); idx++ {
		cnt, err = l.svcCtx.ChatBalanceModel.FindCountEnterChatUserCntRangeUpdateTime(l.ctx, in.ListenerUid[idx], start, end)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ChatRedis.SetTodayListenerEnterChatUserCnt(l.ctx, in.ListenerUid[idx], cnt)
		if err != nil {
			return nil, err
		}
	}

	return &pb.UpdateTodayEnterChatUserCntResp{}, nil
}
