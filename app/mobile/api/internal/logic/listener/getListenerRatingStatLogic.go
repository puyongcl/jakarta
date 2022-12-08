package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRatingStatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerRatingStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRatingStatLogic {
	return &GetListenerRatingStatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerRatingStatLogic) GetListenerRatingStat(req *types.GetListenerRatingStatReq) (resp *types.GetListenerRatingStatResp, err error) {
	var in pb.GetListenerRatingStatReq
	in.ListenerUid = req.ListenerUid
	rs, err := l.svcCtx.ListenerRpc.GetListenerRatingStat(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerRatingStatResp{CommentTagStat: make([]*types.CommentTagPair, 0)}
	_ = copier.Copy(resp, rs)
	return
}
