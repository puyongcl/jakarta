package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStoryReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelStoryReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStoryReplyLogic {
	return &DelStoryReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelStoryReplyLogic) DelStoryReply(req *types.DelStoryReplyReq) (resp *types.DelStoryReplyResp, err error) {
	var in pb.DelStoryReplyReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.BbsRpc.DelStoryReply(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.DelStoryReplyResp{}
	return
}
