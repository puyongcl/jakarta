package stat

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/statistic/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserChannelListLogic {
	return &GetUserChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserChannelListLogic) GetUserChannelList(req *types.GetUserChannelListReq) (resp *types.GetUserChannelListResp, err error) {
	var in pb.GetUserChannelListReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.StatRpc.GetUserChannelList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetUserChannelListResp{
		Channel: make([]string, 0),
	}
	_ = copier.Copy(resp, rs)
	return
}
