package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListByAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListByAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListByAdminLogic {
	return &GetOrderListByAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理后台获取订单列表
func (l *GetOrderListByAdminLogic) GetOrderListByAdmin(in *pb.GetOrderListByAdminReq) (*pb.GetOrderListByAdminResp, error) {
	if in.OrderId != "" {
		data, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
		if err != nil && err != orderPgModel.ErrNotFound {
			return nil, err
		}
		var val pb.RefundChatOrder
		resp := pb.GetOrderListByAdminResp{Sum: 0, List: []*pb.RefundChatOrder{&val}}
		if data != nil {
			l.getRefundOrder(data, &val)
			resp.Sum = 1
		}
		return &resp, nil
	}
	//if len(in.State) <= 0 {
	//	in.State = orderkey.ChatOrderRefundState
	//}
	cnt, err := l.svcCtx.ChatOrderModel.FindCount(l.ctx, in.Uid, in.ListenerUid, in.OrderType, in.State, 0)
	if err != nil {
		return nil, err
	}
	if cnt <= 0 {
		return &pb.GetOrderListByAdminResp{}, nil
	}
	rsp, err := l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, in.OrderType, in.State, orderkey.OrderListTypeAll, "create_time DESC", in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		return &pb.GetOrderListByAdminResp{}, nil
	}
	resp := pb.GetOrderListByAdminResp{Sum: cnt, List: make([]*pb.RefundChatOrder, 0)}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.RefundChatOrder
		l.getRefundOrder(rsp[idx], &val)

		resp.List = append(resp.List, &val)
	}
	return &resp, nil
}

func (l *GetOrderListByAdminLogic) getRefundOrder(data *orderPgModel.ChatOrder, val *pb.RefundChatOrder) {
	_ = copier.Copy(&val, data)
	val.CreateTime = data.CreateTime.Format(db.DateTimeFormat)
	if data.ExpiryTime.Valid {
		val.ExpiryTime = data.ExpiryTime.Time.Format(db.DateTimeFormat)
	}
	if data.StartTime.Valid {
		val.StartTime = data.StartTime.Time.Format(db.DateTimeFormat)
	}
	if data.EndTime.Valid {
		val.EndTime = data.EndTime.Time.Format(db.DateTimeFormat)
	}
	if data.CommentTime.Valid {
		val.CommentTime = data.CommentTime.Time.Format(db.DateTimeFormat)
	}
	if data.FeedbackTime.Valid {
		val.FeedbackTime = data.FeedbackTime.Time.Format(db.DateTimeFormat)
	}
	if data.ApplyRefundTime.Valid {
		val.ApplyRefundTime = data.ApplyRefundTime.Time.Format(db.DateTimeFormat)
	}
}
