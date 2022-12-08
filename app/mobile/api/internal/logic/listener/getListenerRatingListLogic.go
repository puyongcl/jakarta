package listener

import (
	"context"
	"github.com/jinzhu/copier"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRatingListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerRatingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRatingListLogic {
	return &GetListenerRatingListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerRatingListLogic) GetListenerRatingList(req *types.GetListenerRatingListReq) (resp *types.GetListenerRatingListResp, err error) {
	var in pbOrder.GetListenerCommentListReq
	_ = copier.Copy(&in, req)
	uid := ctxdata.GetUidFromCtx(l.ctx)
	in.Uid = uid
	rsp, err := l.svcCtx.OrderRpc.GetListenerCommentList(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerRatingListResp{
		List: make([]*types.ListenerOrderOpinion, 0),
	}
	if len(rsp.List) <= 0 {
		return
	}
	_ = copier.Copy(resp, rsp)
	return
}
