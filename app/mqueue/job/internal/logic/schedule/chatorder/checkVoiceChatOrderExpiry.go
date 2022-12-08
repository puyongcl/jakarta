package chatorder

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"strconv"
	"time"
)

type CheckVoiceChatOrderExpiryHandler struct {
	svcCtx *svc.ServiceContext
}

func NewCheckVoiceChatOrderExpiryHandler(svcCtx *svc.ServiceContext) *CheckVoiceChatOrderExpiryHandler {
	return &CheckVoiceChatOrderExpiryHandler{
		svcCtx: svcCtx,
	}
}

// 1天到期
func (l *CheckVoiceChatOrderExpiryHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取1天后过期的订单
	var pageNo int64 = 1
	now := time.Now()
	endExpiryTime := now.Add(notify.ScheduleNotifyWillExpiryOrderHour * time.Hour).Format(db.DateTimeFormat)
	startExpiryTime := now.Add(notify.ScheduleNotifyWillExpiryOrderHour * time.Hour).Add(-(orderkey.CheckOrderExpiryIntervalMinutes + 1) * time.Minute).Format(db.DateTimeFormat)
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
			logx.WithContext(ctx).Errorf("CheckVoiceChatOrderExpiryHandler GetExpireVoiceChatOrder err:%+v", err)
			return err
		}

		if len(rsp.List) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("CheckVoiceChatOrderExpiryHandler exit.process order cnt:%d", cnt)
			}
			break
		}

		err = l.notify(ctx, rsp.List)
		if err != nil {
			logx.WithContext(ctx).Errorf("CheckVoiceChatOrderExpiryHandler notify err:%+v", err)
			return err
		}
		cnt += len(rsp.List)
	}
	logx.WithContext(ctx).Infof("CheckVoiceChatOrderExpiryHandler done")
	return nil
}

func (l *CheckVoiceChatOrderExpiryHandler) notify(ctx context.Context, inList []*pbOrder.ExpireVoiceChatOrder) (err error) {
	for idx := 0; idx < len(inList); idx++ {
		// 发送提前1天过期消息
		msg := &kqueue.SendImDefineMessage{
			FromUid:           notify.TimOrderNotifyUid,
			ToUid:             inList[idx].Uid,
			MsgType:           notify.DefineNotifyMsgTypeOrderMsg3,
			Title:             notify.DefineNotifyMsgTemplateOrderMsgTitle3,
			Text:              notify.DefineNotifyMsgTemplateOrderMsg3,
			Val1:              strconv.FormatInt(inList[idx].ListenerUid, 10),
			Val2:              inList[idx].OrderId,
			Sync:              tim.TimMsgSyncFromNo,
			RepeatMsgCheckId:  inList[idx].OrderId,
			RepeatMsgCheckSec: orderkey.CheckOrderExpiryIntervalMinutes * 60 * 10,
		}
		var buf []byte
		buf, err = json.Marshal(msg)
		if err != nil {
			return err
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			return err
		}
	}
	return
}
