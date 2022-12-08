package listener

import (
	"context"
	pbOrder "jakarta/app/order/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyCommentLogic {
	return &ReplyCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyCommentLogic) ReplyComment(req *types.ReplyCommentReq) (resp *types.ReplyCommentResp, err error) {
	_, err = l.svcCtx.OrderRpc.ReplyOrderComment(l.ctx, &pbOrder.ReplyOrderCommentReq{
		OrderId:     req.OrderId,
		ListenerUid: req.ListenerUid,
		Reply:       req.Reply,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ReplyCommentResp{}
	return
}
