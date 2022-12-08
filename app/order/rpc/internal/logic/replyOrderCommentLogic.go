package logic

import (
	"context"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/orderkey"
	"jakarta/common/xerr"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyOrderCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyOrderCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyOrderCommentLogic {
	return &ReplyOrderCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  XXX回复用户的订单评价
func (l *ReplyOrderCommentLogic) ReplyOrderComment(in *pb.ReplyOrderCommentReq) (*pb.ReplyOrderCommentResp, error) {
	data, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	if data.Feedback != "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "已经回复，当前无法再次回复")
	}
	newData := new(orderPgModel.ChatOrder)
	newData.Reply = in.Reply
	newData.OrderId = in.OrderId
	err = l.svcCtx.ChatOrderModel.UpdateOrderOpinionAndState(l.ctx, in.OrderId, newData, orderkey.NeedFeedbackOderState)
	if err != nil {
		return nil, err
	}

	return &pb.ReplyOrderCommentResp{}, nil
}
