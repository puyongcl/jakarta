package chatorder

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"time"
)

type UpdateVoiceChatOrderExpiryHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUpdateVoiceChatOrderExpiryHandler(svcCtx *svc.ServiceContext) *UpdateVoiceChatOrderExpiryHandler {
	return &UpdateVoiceChatOrderExpiryHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *UpdateVoiceChatOrderExpiryHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取过期的订单
	var pageNo int64 = 1
	now := time.Now()
	endExpiryTime := now.Format(db.DateTimeFormat)
	startExpiryTime := now.Add(-orderkey.UpdateOrderExpiryIntervalMinutes * time.Minute * 2).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbOrder.GetExpireVoiceChatOrderResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.OrderRpc.GetExpireVoiceChatOrder(ctx, &pbOrder.GetExpireVoiceChatOrderReq{
			PageNo:          pageNo,
			PageSize:        10,
			StartExpiryTime: startExpiryTime,
			EndExpiryTime:   endExpiryTime,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("UpdateVoiceChatOrderExpiryHandler GetExpireVoiceChatOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("UpdateVoiceChatOrderExpiryHandler exit cnt:%d", cnt)
			}
			break
		}

		err = l.update(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("UpdateVoiceChatOrderExpiryHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("UpdateVoiceChatOrderExpiryHandler done")
	return nil
}

func (l *UpdateVoiceChatOrderExpiryHandler) update(ctx context.Context, inList []*pbOrder.ExpireVoiceChatOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		// 更新order state
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(ctx, &pbOrder.DoChatOrderActionReq{
			OrderId:     inList[idx].OrderId,
			Action:      orderkey.ChatOrderStateExpire17,
			OperatorUid: orderkey.DefaultSystemOperatorUid,
		})
		if err != nil {
			return err
		}
	}
	return
}
