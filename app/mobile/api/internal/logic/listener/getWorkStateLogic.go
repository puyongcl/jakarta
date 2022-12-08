package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkStateLogic {
	return &GetWorkStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkStateLogic) GetWorkState(req *types.GetWorkStateReq) (resp *types.GetWorkStateResp, err error) {
	rs, err := l.svcCtx.ListenerRpc.GetWorkState(l.ctx, &pb.GetWorkStateReq{ListenerUid: req.ListenerUid})
	if err != nil {
		return nil, err
	}
	resp = &types.GetWorkStateResp{}
	_ = copier.Copy(resp, rs)
	return
}
