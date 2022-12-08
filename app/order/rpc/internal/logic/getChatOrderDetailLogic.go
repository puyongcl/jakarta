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

type GetChatOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatOrderDetailLogic {
	return &GetChatOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  系统内部获取订单详情
func (l *GetChatOrderDetailLogic) GetChatOrderDetail(in *pb.GetChatOrderDetailReq) (*pb.GetChatOrderDetailResp, error) {
	rsp, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil && err != orderPgModel.ErrNotFound {
		return nil, err
	}
	resp := &pb.GetChatOrderDetailResp{Order: &pb.ShortChatOrder{}}
	if rsp != nil {
		_ = copier.Copy(resp.Order, rsp)
		resp.Order.CreateTime = rsp.CreateTime.Format(db.DateTimeFormat)
		if rsp.ExpiryTime.Valid {
			resp.Order.ExpiryTime = rsp.ExpiryTime.Time.Format(db.DateTimeFormat)
		}
		if rsp.StartTime.Valid {
			resp.Order.StartTime = rsp.StartTime.Time.Format(db.DateTimeFormat)
		}
		if rsp.EndTime.Valid {
			resp.Order.EndTime = rsp.EndTime.Time.Format(db.DateTimeFormat)
		}
		if rsp.CommentTime.Valid {
			resp.Order.CommentTime = rsp.CommentTime.Time.Format(db.DateTimeFormat)
		}
		if rsp.FeedbackTime.Valid {
			resp.Order.FeedbackTime = rsp.FeedbackTime.Time.Format(db.DateTimeFormat)
		}
		if rsp.ReplyTime.Valid {
			resp.Order.ReplyTime = rsp.ReplyTime.Time.Format(db.DateTimeFormat)
		}
	}

	return resp, nil
}
