package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryListByOwnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStoryListByOwnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryListByOwnLogic {
	return &GetStoryListByOwnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStoryListByOwnLogic) GetStoryListByOwn(req *types.GetStoryListByOwnReq) (resp *types.GetStoryListByOwnResp, err error) {
	var in pb.GetStoryListByOwnReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.GetStoryListByOwn(l.ctx, &in)
	if err != nil {
		return nil, err
	}

	resp = &types.GetStoryListByOwnResp{List: make([]*types.Story, 0)}
	_ = copier.Copy(&resp, rsp)
	return
}
