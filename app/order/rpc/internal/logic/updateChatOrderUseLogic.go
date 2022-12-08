package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/xerr"
	"strconv"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatOrderUseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatOrderUseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatOrderUseLogic {
	return &UpdateChatOrderUseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新订单的使用情况
func (l *UpdateChatOrderUseLogic) UpdateChatOrderUse(in *pb.UpdateChatOrderUseReq) (*pb.UpdateChatOrderUseResp, error) {
	if in.Uid == 0 || in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	var err error
	var vol []*orderPgModel.ChatOrder
	// 查询当前可用订单
	vol, err = l.svcCtx.ChatOrderModel.Find(l.ctx, in.Uid, in.ListenerUid, in.OrderType, orderkey.CanStopOrderState, orderkey.OrderListTypeAll, "create_time ASC", 1, 100)
	if err != nil {
		return nil, err
	}
	if len(vol) <= 0 { //
		return nil, xerr.NewGrpcErrCodeMsg(xerr.OrderError, "没有找到可以结束的订单")
	}

	switch in.Action {
	case orderkey.ChatOrderStateUsing4, orderkey.ChatOrderStateUseOutWaitUserConfirm13: // 正常使用
		switch in.OrderType {
		case orderkey.ListenerOrderTypeTextChat:
			err = l.UpdateTextChatOrderUse(in, vol)
		case orderkey.ListenerOrderTypeVoiceChat:
			err = l.UpdateVoiceChatOrderUse(in, vol)
		default:
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "error order type")
		}
	case orderkey.ChatOrderStateUserStopService12: // 用户主动结束服务
		switch in.OrderType {
		case orderkey.ListenerOrderTypeTextChat:
			err = l.updateTextChatStop(in, vol)
		case orderkey.ListenerOrderTypeVoiceChat:
			err = l.updateVoiceChatStop(in, vol)
		default:
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "error order type")
		}
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "error order action")
	}
	if err != nil {
		return nil, err
	}

	// 发送通知
	var orderId string
	if len(vol) == 1 {
		orderId = vol[0].OrderId
	}

	var amount int64
	for k, _ := range vol {
		amount += vol[k].ActualAmount
	}
	l.notify(in, orderId, amount)
	return &pb.UpdateChatOrderUseResp{}, nil
}

//  更新订单的使用情况
func (l *UpdateChatOrderUseLogic) UpdateVoiceChatOrderUse(in *pb.UpdateChatOrderUseReq, vol []*orderPgModel.ChatOrder) error {
	// 更新使用时间
	// 剩余待分配的时间
	var err error
	currChatUsedMin := in.UsedMinute
	for idx := 0; idx < len(vol); idx++ {
		if currChatUsedMin <= 0 {
			break
		}
		// 更新开始时间
		var startTime *time.Time
		// 当前订单使用时间
		var currOrderUsedMin int64
		var endTime *time.Time
		state := in.Action
		if vol[idx].UsedChatMinute == 0 {
			var tt time.Time
			tt, err = time.ParseInLocation(db.DateTimeFormat, in.StartTime, time.Local)
			if err != nil {
				return err
			}
			startTime = &tt
			state = orderkey.ChatOrderStateUsing4
		}
		// 未使用时间
		notUsedMin := (vol[idx].ChatUnitMinute * vol[idx].BuyUnit) - vol[idx].UsedChatMinute
		if notUsedMin > 0 {
			if currChatUsedMin < notUsedMin { // 部分使用
				currOrderUsedMin = vol[idx].UsedChatMinute + currChatUsedMin
				currChatUsedMin = 0
				state = orderkey.ChatOrderStateUsing4

			} else { // 刚好用完或者不足
				currOrderUsedMin = vol[idx].UsedChatMinute + notUsedMin
				currChatUsedMin -= notUsedMin
				var tt time.Time
				tt, err = time.ParseInLocation(db.DateTimeFormat, in.StopTime, time.Local)
				if err != nil {
					return err
				}
				endTime = &tt
				state = orderkey.ChatOrderStateUseOutWaitUserConfirm13
			}
		} else {
			var tt time.Time
			tt, err = time.ParseInLocation(db.DateTimeFormat, in.StopTime, time.Local)
			if err != nil {
				return err
			}
			endTime = &tt
			state = orderkey.ChatOrderStateUseOutWaitUserConfirm13
		}

		err = l.updateOrder(vol[idx], state, currOrderUsedMin, startTime, endTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// 用户主动停止服务
func (l *UpdateChatOrderUseLogic) updateVoiceChatStop(in *pb.UpdateChatOrderUseReq, vol []*orderPgModel.ChatOrder) error {
	now := time.Now()
	var err error
	var startTime, endTime *time.Time
	var usedMin int64
	for idx := 0; idx < len(vol); idx++ {
		// 使用时间计算
		switch vol[idx].OrderState {
		case orderkey.ChatOrderStatePaySuccess3:
			startTime = nil
			endTime = nil
			usedMin = 0
		case orderkey.ChatOrderStateUsing4:
			if vol[idx].StartTime.Valid {
				usedMin = int64(now.Sub(vol[idx].StartTime.Time).Minutes())
			}
			endTime = &now

		default:
			continue
		}
		err = l.updateOrder(vol[idx], in.Action, usedMin, startTime, endTime)
		if err != nil {
			return err
		}
	}
	return nil
}

// 用户主动停止服务
func (l *UpdateChatOrderUseLogic) updateTextChatStop(in *pb.UpdateChatOrderUseReq, vol []*orderPgModel.ChatOrder) error {
	now := time.Now()
	var startTime, endTime *time.Time
	var usedMin int64
	var err error
	for idx := 0; idx < len(vol); idx++ {
		if vol[idx].ExpiryTime.Valid && vol[idx].ExpiryTime.Time.After(now) {
			if vol[idx].StartTime.Valid {
				usedMin = int64(now.Sub(vol[idx].StartTime.Time).Minutes() + 1)
			}
			endTime = &now
		} else {
			usedMin = int64(now.Sub(vol[idx].CreateTime).Minutes())
			startTime = &(vol[idx].CreateTime)
			endTime = &now
		}
		err = l.updateOrder(vol[idx], in.Action, usedMin, startTime, endTime)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *UpdateChatOrderUseLogic) UpdateTextChatOrderUse(in *pb.UpdateChatOrderUseReq, vol []*orderPgModel.ChatOrder) error {
	now := time.Now()
	// 更新订单
	var err error
	for idx := 0; idx < len(vol); idx++ {
		if vol[idx].ExpiryTime.Valid && vol[idx].ExpiryTime.Time.Before(now.Add(time.Minute)) { // 有效期到了
			err = l.updateOrder(vol[idx], orderkey.ChatOrderStateUseOutWaitUserConfirm13, vol[idx].BuyUnit*vol[idx].ChatUnitMinute, nil, &now)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *UpdateChatOrderUseLogic) updateOrder(orderData *orderPgModel.ChatOrder, state int64, usedMin int64, startTime, endTime *time.Time) error {
	ndc := NewDoChatOrderActionLogic(l.ctx, l.svcCtx)
	in2 := pb.DoChatOrderActionReq{
		OrderId:     orderData.OrderId,
		OperatorUid: orderData.Uid,
		Action:      state,
		OrderType:   orderData.OrderType,
		ListenerUid: orderData.ListenerUid,
	}
	_, err := ndc.OneOrderAction(usedMin, startTime, endTime, &in2, orderData)
	if err != nil {
		return err
	}
	return nil
}

func (l *UpdateChatOrderUseLogic) notify(in *pb.UpdateChatOrderUseReq, orderId string, amount int64) {
	// 发送系统通知
	var msgs []*kqueue.SendImDefineMessage
	switch in.Action {
	case orderkey.ChatOrderStateUsing4:
		msg := kqueue.SendImDefineMessage{
			FromUid: in.ListenerUid,
			ToUid:   in.Uid,
			MsgType: notify.DefineNotifyMsgTypeChatMsg5,
			Title:   "",
			Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateChatMsg5, in.UsedMinute),
			Val1:    strconv.FormatInt(in.UsedMinute, 10),
			Val2:    orderId,
			Val3:    strconv.FormatInt(in.Uid, 10),
			Val4:    strconv.FormatInt(in.ListenerUid, 10),
			Sync:    tim.TimMsgSyncFromYes,
		}
		msgs = append(msgs, &msg)
	case orderkey.ChatOrderStateUseOutWaitUserConfirm13, orderkey.ChatOrderStateUserStopService12: // 用完或主动结束
		// 发送服务结束
		// 对用户
		msg1 := kqueue.SendImDefineMessage{
			FromUid: in.ListenerUid,
			ToUid:   in.Uid,
			MsgType: notify.DefineNotifyMsgTypeChatMsg22,
			Title:   "",
			Text:    notify.DefineNotifyMsgTemplateChatMsg22,
			Val1:    "特邀请你对服务评价",
			Val2:    orderId,
			Val3:    strconv.FormatInt(in.Uid, 10),
			Val4:    strconv.FormatInt(in.ListenerUid, 10),
			Val5:    strconv.FormatInt(amount, 10),
			Sync:    tim.TimMsgSyncFromNo,
		}
		// 对XXX
		msg2 := kqueue.SendImDefineMessage{
			FromUid: in.Uid,
			ToUid:   in.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeChatMsg23,
			Title:   "",
			Text:    notify.DefineNotifyMsgTemplateChatMsg23,
			Val1:    "请你对XX者反馈鼓励",
			Val2:    orderId,
			Val3:    strconv.FormatInt(in.Uid, 10),
			Val4:    strconv.FormatInt(in.ListenerUid, 10),
			Val5:    strconv.FormatInt(amount, 10),
			Sync:    tim.TimMsgSyncFromNo,
		}

		if in.OrderType == orderkey.ListenerOrderTypeVoiceChat {
			msg1.MsgType = notify.DefineNotifyMsgTypeChatMsg29
			msg1.Text = notify.DefineNotifyMsgTemplateChatMsg29
			msg2.MsgType = notify.DefineNotifyMsgTypeChatMsg30
			msg2.Text = notify.DefineNotifyMsgTemplateChatMsg30
		}
		msgs = append(msgs, &msg1, &msg2)
	}

	if len(msgs) <= 0 {
		return
	}

	var buf []byte
	var err error
	for k, _ := range msgs {
		buf, err = json.Marshal(msgs[k])
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UpdateChatOrderUseLogic json marshal err:%+v", err)
			return
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UpdateChatOrderUseLogic Push err:%+v", err)
			return
		}
	}
}
