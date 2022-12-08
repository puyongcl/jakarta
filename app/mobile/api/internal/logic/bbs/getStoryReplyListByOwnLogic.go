package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryReplyListByOwnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStoryReplyListByOwnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryReplyListByOwnLogic {
	return &GetStoryReplyListByOwnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStoryReplyListByOwnLogic) GetStoryReplyListByOwn(req *types.GetStoryReplyListByOwnReq) (resp *types.GetStoryReplyListByOwnResp, err error) {
	var in pb.GetStoryReplyListByOwnReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.GetStoryReplyListByOwn(l.ctx, &in)
	if err != nil {
		return nil, err
	}

	resp = &types.GetStoryReplyListByOwnResp{List: make([]*types.StoryReply, 0)}
	_ = copier.Copy(&resp, rsp)
	return
}
