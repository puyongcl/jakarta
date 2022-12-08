package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelStoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStoryLogic {
	return &DelStoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelStoryLogic) DelStory(req *types.DelStoryReq) (resp *types.DelStoryResp, err error) {
	var in pb.DelStoryReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.BbsRpc.DelStory(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.DelStoryResp{}
	return
}
