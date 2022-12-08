package stat

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/statistic/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLTVStatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLTVStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLTVStatLogic {
	return &GetLTVStatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLTVStatLogic) GetLTVStat(req *types.GetLifeTimeValueStatReq) (resp *types.GetLifeTimeValueStatResp, err error) {
	var in pb.GetLifeTimeValueStatReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.StatRpc.GetLifeTimeValueStat(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetLifeTimeValueStatResp{
		List: make([]*types.LifeTimeValueStat, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
