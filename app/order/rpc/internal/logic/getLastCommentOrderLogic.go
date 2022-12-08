package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastCommentOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastCommentOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastCommentOrderLogic {
	return &GetLastCommentOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取最近一条评价
func (l *GetLastCommentOrderLogic) GetLastCommentOrder(in *pb.GetLastCommentOrderReq) (*pb.GetLastCommentOrderResp, error) {
	rsp, err := l.svcCtx.ChatOrderModel.FindLastCommentOrder(l.ctx, in.ListenerUid, in.Star)
	if err != nil && err != orderPgModel.ErrNotFound {
		return nil, err
	}

	resp := pb.GetLastCommentOrderResp{}
	if rsp != nil {
		_ = copier.Copy(&resp, rsp)
		resp.CommentTime = rsp.CommentTime.Time.Format(db.DateTimeFormat)
	}
	return &resp, nil
}
