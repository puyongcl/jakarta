package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatOrderListLogic {
	return &GetChatOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  订单列表
func (l *GetChatOrderListLogic) GetChatOrderList(in *pb.GetChatOrderListReq) (*pb.GetChatOrderListResp, error) {
	rsp := make([]*orderPgModel.ChatOrder, 0)
	var err error
	switch in.ListType {
	case orderkey.OrderListTypeAll: // 全部
		rsp, err = l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, 0, []int64{}, orderkey.OrderListTypeAll, "create_time DESC", in.PageNo, in.PageSize)
	case orderkey.OrderListTypeNeedProcess: // 未完成
		// 服务中
		var rsp1 []*orderPgModel.ChatOrder
		rsp1, err = l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, 0, []int64{orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateUsing4, orderkey.ChatOrderStateTextOrderPaySuccess24}, orderkey.OrderListTypeNeedProcess, "create_time DESC", in.PageNo, in.PageSize)
		if err != nil {
			return nil, err
		}
		// 未反馈
		var rsp2 []*orderPgModel.ChatOrder
		rsp2, err = l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, 0, []int64{}, orderkey.OrderListTypeNeedFeedback, "create_time DESC", in.PageNo, in.PageSize)
		if err != nil {
			return nil, err
		}

		// 待处理退款
		var rsp3 []*orderPgModel.ChatOrder
		rsp3, err = l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, 0, []int64{orderkey.ChatOrderStateApplyRefund5}, orderkey.OrderListTypeNeedAgreeRefund, "create_time DESC", in.PageNo, in.PageSize)
		if err != nil {
			return nil, err
		}
		if len(rsp3) > 0 {
			rsp = append(rsp, rsp3...)
		}
		if len(rsp1) > 0 {
			rsp = append(rsp, rsp1...)
		}
		if len(rsp2) > 0 {
			rsp = append(rsp, rsp2...)
		}
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "error list type")
	}

	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		return &pb.GetChatOrderListResp{}, nil
	}
	resp := pb.GetChatOrderListResp{List: make([]*pb.ShortChatOrder, 0)}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.ShortChatOrder
		_ = copier.Copy(&val, rsp[idx])
		val.CreateTime = rsp[idx].CreateTime.Format(db.DateTimeFormat)
		if rsp[idx].ExpiryTime.Valid {
			val.ExpiryTime = rsp[idx].ExpiryTime.Time.Format(db.DateTimeFormat)
		}
		if rsp[idx].StartTime.Valid {
			val.StartTime = rsp[idx].StartTime.Time.Format(db.DateTimeFormat)
		}
		if rsp[idx].EndTime.Valid {
			val.EndTime = rsp[idx].EndTime.Time.Format(db.DateTimeFormat)
		}
		if rsp[idx].CommentTime.Valid {
			val.CommentTime = rsp[idx].CommentTime.Time.Format(db.DateTimeFormat)
		}
		if rsp[idx].FeedbackTime.Valid {
			val.FeedbackTime = rsp[idx].FeedbackTime.Time.Format(db.DateTimeFormat)
		}
		if rsp[idx].ReplyTime.Valid {
			val.ReplyTime = rsp[idx].ReplyTime.Time.Format(db.DateTimeFormat)
		}
		// 订单状态显示
		val.StatusMark = getOrderStatusMark(val.OrderState, rsp[idx].ExpiryTime.Time, rsp[idx].UpdateTime)
		resp.List = append(resp.List, &val)
	}
	return &resp, nil
}

func getOrderStatusMark(orderState int64, expiryTime time.Time, updateTime time.Time) string {
	switch orderState { // TODO 退款相关状态 不做提示
	case orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateUsing4:
		td := expiryTime.Sub(time.Now())
		m := int(td.Minutes())
		h := int(td.Hours())
		return fmt.Sprintf("有效期还剩%d小时%d分钟", h, m-(h*60))
	case orderkey.ChatOrderStateUserStopService12, orderkey.ChatOrderStateUseOutWaitUserConfirm13:
		td := updateTime.AddDate(0, 0, orderkey.AutoGoodCommentAfterStopDay).Sub(time.Now())
		h := int(td.Hours())
		d := h / 24
		if d > 0 {
			return fmt.Sprintf("%d天后自动好评和确认完成", d)
		}
		return fmt.Sprintf("%d小时后自动好评和确认完成", h)

	case orderkey.ChatOrderStateNot5StarWaitConfirm15: // 非好评后自动确认完成
		td := updateTime.AddDate(0, 0, orderkey.AutoFinishAfterNotGoodCommentDay).Sub(time.Now())
		h := int(td.Hours())
		d := h / 24
		if d > 0 {
			return fmt.Sprintf("%d天后自动确认完成", d)
		}
		return fmt.Sprintf("%d小时后自动确认完成", h)

	case orderkey.ChatOrderStateSecondApplyRefund10: // 申请客服介入退款后 不处理 自动处理为 确认完成
		//td := updateTime.AddDate(0, 0, orderkey.AutoFinishAfterSecondApplyRefundDay).Sub(time.Now())
		//h := int(td.Hours())
		//d := h / 24
		//if d > 0 {
		//	return fmt.Sprintf("%d天后自动确认完成", d)
		//}
		//return fmt.Sprintf("%d小时后自动确认完成", h)
	case orderkey.ChatOrderStateListenerRefuseRefund9: // XXX拒绝退款 自动完成
		//td := updateTime.AddDate(0, 0, orderkey.AutoFinishAfterListenerRefuseRefund).Sub(time.Now())
		//h := int(td.Hours())
		//d := h / 24
		//if d > 0 {
		//	return fmt.Sprintf("%d天后自动确认完成", d)
		//}
		//return fmt.Sprintf("%d小时后自动确认完成", h)
	case orderkey.ChatOrderStateApplyRefund5: // 申请退款 几天内不处理自动同意
	case orderkey.ChatOrderStateListenerAgreeRefund7, orderkey.ChatOrderStateAdminAgreeRefund11,
		orderkey.ChatOrderStateAutoAgreeRefund6: // 同意退款多久到账
	default:

	}
	return ""
}
