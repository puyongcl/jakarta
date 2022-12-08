package stat

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/statistic/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDailyStatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDailyStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDailyStatLogic {
	return &GetDailyStatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDailyStatLogic) GetDailyStat(req *types.GetDailyStatListReq) (resp *types.GetDailyStatListResp, err error) {
	var in pb.GetDailyStatListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.StatRpc.GetDailyStatList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetDailyStatListResp{
		List: make([]*types.DailyStat, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
