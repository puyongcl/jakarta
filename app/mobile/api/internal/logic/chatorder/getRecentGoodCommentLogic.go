package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentGoodCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecentGoodCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecentGoodCommentLogic {
	return &GetRecentGoodCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecentGoodCommentLogic) GetRecentGoodComment(req *types.GetRecentGoodCommentReq) (resp *types.GetRecentGoodCommentResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	//
	var in pbListener.GetRecommendListenerReq
	in.Uid = uid
	in.PageNo = req.PageNo
	in.PageSize = req.PageSize

	rsp, err := l.svcCtx.ListenerRpc.GetRecommendListenerList(l.ctx, &in)
	if err != nil {
		return
	}
	if len(rsp.Listener) <= 0 {
		resp = &types.GetRecentGoodCommentResp{List: make([]*types.RecentGoodComment, 0)}
		return
	}

	var in2 pbOrder.GetListenerRecentGoodCommentReq
	in2.Listener = make([]*pbOrder.ShortListenerProfile, 0)
	for idx := 0; idx < len(rsp.Listener); idx++ {
		var val pbOrder.ShortListenerProfile
		val.ListenerUid = rsp.Listener[idx].ListenerUid
		val.ListenerNickName = rsp.Listener[idx].ListenerNickName
		val.ListenerAvatar = rsp.Listener[idx].ListenerAvatar

		in2.Listener = append(in2.Listener, &val)
	}
	//
	rsp2, err := l.svcCtx.OrderRpc.GetListenerRecentGoodComment(l.ctx, &in2)
	if err != nil {
		return nil, err
	}
	resp = &types.GetRecentGoodCommentResp{List: make([]*types.RecentGoodComment, 0)}
	_ = copier.Copy(resp, rsp2)
	return
}
