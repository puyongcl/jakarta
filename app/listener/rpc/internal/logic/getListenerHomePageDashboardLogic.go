package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/dashboard"
	"jakarta/common/key/listenerkey"
)

type GetListenerHomePageDashboardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerHomePageDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerHomePageDashboardLogic {
	return &GetListenerHomePageDashboardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX首页统计数据
func (l *GetListenerHomePageDashboardLogic) GetListenerHomePageDashboard(in *pb.GetListenerHomePageDashboardReq) (*pb.GetListenerHomePageDashboardResp, error) {
	data, err := l.svcCtx.StatRedis.GetListenerDashboard(l.ctx, in.ListenerUid)
	if err != nil {
		return nil, err
	}

	resp := pb.GetListenerHomePageDashboardResp{}

	if data != nil {
		var rsp dashboard.ListenerDashboard
		err = dashboard.TransferListenerDashboardData(data, &rsp)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&resp, &rsp)
	}

	// 建议
	var stat *listenerPgModel.ListenerDashboardStat
	stat, err = l.svcCtx.ListenerDashboardStatModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil {
		return nil, err
	}
	for idx := 0; idx < len(stat.Suggestion); idx++ {
		resp.Suggestion = append(resp.Suggestion, listenerkey.Suggestion[stat.Suggestion[idx]])
	}

	return &resp, nil
}
