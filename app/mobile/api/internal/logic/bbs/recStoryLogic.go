package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecStoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecStoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecStoryLogic {
	return &RecStoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecStoryLogic) RecStory(req *types.GetRecStoryListByUserReq) (resp *types.GetRecStoryListByUserResp, err error) {
	var in pb.GetRecStoryListByUserReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.GetRecStoryListByUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}

	resp = &types.GetRecStoryListByUserResp{List: make([]*types.Story, 0)}
	_ = copier.Copy(&resp, rsp)
	return
}
