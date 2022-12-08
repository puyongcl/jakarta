package chatorder

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"time"
)

type AutoCommentChatOrderHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoCommentChatOrderHandler(svcCtx *svc.ServiceContext) *AutoCommentChatOrderHandler {
	return &AutoCommentChatOrderHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoCommentChatOrderHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取到时的订单
	var pageNo int64 = 1
	bTime := time.Now().AddDate(0, 0, -orderkey.AutoGoodCommentAfterStopDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbOrder.GetAutoProcessOrderResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.OrderRpc.GetAutoProcessOrder(ctx, &pbOrder.GetAutoProcessOrderReq{
			PageNo:     pageNo,
			PageSize:   10,
			BeforeTime: bTime,
			State:      orderkey.CanAutoGoodCommentAndFinishOrderState,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoCommentChatOrderHandler GetAutoProcessOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoCommentChatOrderHandler exit. process order cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoCommentChatOrderHandler update err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("AutoCommentChatOrderHandler done")
	return nil
}

func (l *AutoCommentChatOrderHandler) update(ctx context.Context, inList []*pbOrder.AutoProcessOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &pbOrder.DoChatOrderActionReq{
			OrderId:     inList[idx].OrderId,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
			Action:      orderkey.ChatOrderStateAutoCommentFinish18,
			Star:        listenerkey.Rating5Star,
			Tag:         []int64{},
			SendMsg:     db.Disable,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
