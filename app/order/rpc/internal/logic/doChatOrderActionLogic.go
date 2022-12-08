package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoChatOrderActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDoChatOrderActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoChatOrderActionLogic {
	return &DoChatOrderActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  订单操作
func (l *DoChatOrderActionLogic) DoChatOrderAction(in *pb.DoChatOrderActionReq) (*pb.DoChatOrderActionResp, error) {
	orderData, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil && err != orderPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, err.Error())
	}
	if orderData == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.OrderError, "未查询到订单")
	}
	// 计算当前订单 停止的时间
	if in.Action == orderkey.ChatOrderStateUserStopService12 {
		usedMin, startTime, endTime := l.stopOrder(orderData)
		return l.OneOrderAction(usedMin, startTime, endTime, in, orderData)
	}

	return l.OneOrderAction(0, nil, nil, in, orderData)
}

func (l *DoChatOrderActionLogic) stopOrder(orderData *orderPgModel.ChatOrder) (usedMin int64, startTime, endTime *time.Time) {
	now := time.Now()

	switch orderData.OrderState {
	case orderkey.ChatOrderStatePaySuccess3:
		startTime = nil
		endTime = nil
		usedMin = 0

	case orderkey.ChatOrderStateTextOrderPaySuccess24:
		if orderData.ExpiryTime.Valid && orderData.ExpiryTime.Time.After(now) {
			if orderData.StartTime.Valid {
				usedMin = int64(now.Sub(orderData.StartTime.Time).Minutes())
			}
			endTime = &now
		} else if orderData.ExpiryTime.Valid {
			endTime = &(orderData.ExpiryTime.Time)
			usedMin = orderData.BuyUnit * orderData.ChatUnitMinute
		} else {
			endTime = &now
			usedMin = orderData.BuyUnit * orderData.ChatUnitMinute
		}

	case orderkey.ChatOrderStateUsing4:
		if orderData.StartTime.Valid {
			usedMin = int64(now.Sub(orderData.StartTime.Time).Minutes())
		}
		endTime = &now

	default:

	}
	return
}

func (l *DoChatOrderActionLogic) OneOrderAction(usedChatMinute int64, startTime, endTime *time.Time, in *pb.DoChatOrderActionReq, orderData *orderPgModel.ChatOrder) (*pb.DoChatOrderActionResp, error) {
	rsp, err := l.action(usedChatMinute, startTime, endTime, in, orderData)
	createActionLog(l.svcCtx.ChatOrderStatusLogModel, l.ctx, in.OrderId, in.Action, in.OperatorUid, in.Remark, err)
	if err != nil {
		return nil, err
	}
	return rsp, err
}

func createActionLog(logModel orderPgModel.ChatOrderStatusLogModel, ctx context.Context, orderId string, action int64, operatorUid int64, remark string, errIn error) {
	logData := orderPgModel.ChatOrderStatusLog{
		Id:          uniqueid.GenDataId(),
		OrderId:     orderId,
		State:       action,
		OperatorUid: operatorUid,
		Remark:      remark,
	}
	if errIn != nil {
		logData.ActionResult = 1
		logData.Remark += fmt.Sprintf("#%+v", errIn)
	}
	_, err := logModel.Insert(ctx, &logData)
	if err != nil {
		logx.WithContext(ctx).Errorf("createActionLog data:%+v err:%+v", logData, err)
		return
	}
	return
}

var errTradeState = xerr.NewGrpcErrCodeMsg(xerr.OrderError, "状态校验错误")

func (l *DoChatOrderActionLogic) action(usedChatMinute int64, startTime, endTime *time.Time, in *pb.DoChatOrderActionReq, orderData *orderPgModel.ChatOrder) (*pb.DoChatOrderActionResp, error) {
	var err error
	resp := &pb.DoChatOrderActionResp{}
	action := in.Action
	// 校验状态
	var states []int64
	err, states = l.verifyState(action, orderData)
	if err != nil {
		return nil, err
	}

	// 更新状态
	err = l.updateOrder(usedChatMinute, startTime, endTime, action, in, orderData, states)
	if err != nil {
		return nil, err
	}

	// 后续操作
	err = l.next(action, in, orderData)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (l *DoChatOrderActionLogic) next(action int64, in *pb.DoChatOrderActionReq, orderData *orderPgModel.ChatOrder) (err error) {
	var addRepeat int64 // 是否增加复购人数
	var addUser int64   // 是否增加服务人数

	switch action {
	case orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateTextOrderPaySuccess24: // 复购人数
		var b bool
		b, err = l.isRepeatUser(orderData.Uid, orderData.ListenerUid)
		if err != nil {
			return
		}
		if b {
			addRepeat = 1
		}
	case orderkey.ChatOrderStateUserStopService12, orderkey.ChatOrderStateUseOutWaitUserConfirm13: // 服务人数
		var b bool
		b, err = l.isNewUser(orderData.Uid, orderData.ListenerUid)
		if err != nil {
			return
		}
		if b && orderData.UsedChatMinute > 0 {
			addUser = 1
		}
	default:

	}

	msg := kqueue.UpdateChatOrderActionMessage{
		OrderId:          orderData.OrderId,
		ListenerUid:      orderData.ListenerUid,
		ListenerNickName: orderData.ListenerNickName,
		Uid:              orderData.Uid,
		NickName:         orderData.NickName,
		Action:           action,
		PaidAmount:       orderData.ActualAmount,
		Star:             orderData.Star,
		BuyMinute:        orderData.BuyUnit * orderData.ChatUnitMinute,
		UsedMinute:       orderData.UsedChatMinute,
		OrderType:        orderData.OrderType,
		OrderCreateTime:  orderData.CreateTime.Format(db.DateTimeFormat),
		Reason:           orderData.RefundReason,
		SendMsg:          in.SendMsg,
		CommentTag:       orderData.CommentTag,
		Comment:          orderData.Comment,
		AddRepeat:        addRepeat,
		AddUser:          addUser,
	}
	if orderData.ExpiryTime.Valid {
		msg.TextExpireTime = orderData.ExpiryTime.Time.Format(db.DateTimeFormat)
	}

	buf, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueUpdateOrderActionClient.Push(string(buf))
	if err != nil {
		return err
	}
	logx.WithContext(l.ctx).Infof("DoChatOrderActionLogic push done order_id:%s time:%s", msg.OrderId, time.Now().Format(db.DateTimeFormat))
	return
}

func (l *DoChatOrderActionLogic) updateOrder(usedChatMinute int64, startTime, endTime *time.Time, action int64, in *pb.DoChatOrderActionReq, orderData *orderPgModel.ChatOrder, states []int64) (err error) {
	switch action {
	case orderkey.ChatOrderState5StartRatingAndFinish14, orderkey.ChatOrderStateNot5StarWaitConfirm15, orderkey.ChatOrderStateAutoCommentFinish18: // 评价
		if orderData.Comment != "" || len(orderData.CommentTag) > 0 {
			return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "已经存在评价，当前无法再次评价")
		}
		newData := new(orderPgModel.ChatOrder)
		newData.Comment = in.Comment
		orderData.Comment = in.Comment
		newData.CommentTag = in.Tag
		orderData.CommentTag = in.Tag
		newData.Star = in.Star
		orderData.Star = in.Star
		newData.OrderId = in.OrderId
		newData.OrderState = action
		err = l.svcCtx.ChatOrderModel.UpdateOrderOpinionAndState(l.ctx, in.OrderId, newData, states)
		if err != nil {
			return err
		}

	case orderkey.ChatOrderStateAutoAgreeRefund6, orderkey.ChatOrderStateApplyRefund5: // 首次申请退款（符合条件自动同意）
		newData := new(orderPgModel.ChatOrder)
		if len(in.Tag) > 1 {
			newData.RefundReasonTag = in.Tag[0]
			orderData.RefundReasonTag = in.Tag[0]
		}

		if in.Remark == "" && len(in.Tag) > 1 {
			in.Remark, _ = orderkey.RefundReasonTag[int(in.Tag[0])]
		}
		if in.Remark != "" {
			newData.RefundReason = in.Remark
			orderData.RefundReason = in.Remark
		}
		if in.Additional != "" {
			newData.Additional = in.Additional
			orderData.Additional = in.Additional
		}
		if in.Attachment != "" {
			newData.Attachment = in.Attachment
			orderData.Attachment = in.Attachment
		}
		newData.ApplyRefundTime.Time = time.Now()
		newData.ApplyRefundTime.Valid = true
		orderData.ApplyRefundTime.Time = time.Now()
		orderData.ApplyRefundTime.Valid = true
		newData.OrderState = action
		err = l.svcCtx.ChatOrderModel.UpdateApplyRefund(l.ctx, in.OrderId, newData, states)
		if err != nil {
			return err
		}

	case orderkey.ChatOrderStateSecondApplyRefund10: // XXX拒绝退款 用户申请客服介入
		newData := new(orderPgModel.ChatOrder)
		if in.Additional != "" {
			newData.Additional = in.Additional
			orderData.Additional = in.Additional
		}
		if in.Attachment != "" {
			newData.Attachment = in.Attachment
			orderData.Attachment = in.Attachment
		}
		newData.OrderState = action
		err = l.svcCtx.ChatOrderModel.UpdateApplyRefund(l.ctx, in.OrderId, newData, states)
		if err != nil {
			return err
		}

	case orderkey.ChatOrderStateListenerRefuseRefund9, orderkey.ChatOrderStateListenerAgreeRefund7, orderkey.ChatOrderStateAdminAgreeRefund11: // 退款
		newData := new(orderPgModel.ChatOrder)
		if in.Remark != "" {
			if orderData.RefundCheckRemark != "" {
				newData.RefundCheckRemark = orderData.RefundCheckRemark + "#" + in.Remark
				orderData.RefundCheckRemark = newData.RefundCheckRemark
			} else {
				newData.RefundCheckRemark = in.Remark
				orderData.RefundCheckRemark = newData.RefundCheckRemark
			}
			if len(newData.RefundCheckRemark) > 510 {
				newData.RefundCheckRemark = newData.RefundCheckRemark[:510]
				orderData.RefundCheckRemark = newData.RefundCheckRemark
			}
		}

		newData.OrderState = action
		err = l.svcCtx.ChatOrderModel.UpdateApplyRefund(l.ctx, in.OrderId, newData, states)
		if err != nil {
			return err
		}

	case orderkey.ChatOrderStateFinishRefund8: // 退款成功
		newData := new(orderPgModel.ChatOrder)
		newData.RefundSuccessTime.Time = time.Now()
		newData.RefundSuccessTime.Valid = true
		orderData.RefundSuccessTime = newData.RefundSuccessTime
		newData.OrderState = action
		err = l.svcCtx.ChatOrderModel.UpdateApplyRefund(l.ctx, in.OrderId, newData, states)
		if err != nil {
			return err
		}

	case orderkey.ChatOrderStateUsing4, orderkey.ChatOrderStateUseOutWaitUserConfirm13, orderkey.ChatOrderStateUserStopService12: // 开始使用\使用完\主动结束
		err = l.svcCtx.ChatOrderModel.UpdateOrderUse(l.ctx, orderData.OrderId, usedChatMinute, startTime, endTime, action, states)
		if err != nil {
			return err
		}
		if usedChatMinute != 0 {
			orderData.UsedChatMinute = usedChatMinute
		}

	case orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateTextOrderPaySuccess24: // 订单支付成功
		// 截止使用日期 文字订单为即时订单，多次下单时间累计，没有截止日期的概念
		now := time.Now()
		var sTime, exTime *time.Time
		switch orderData.OrderType {
		case orderkey.ListenerOrderTypeTextChat:
			// 到期时间设置为整分钟时刻 0秒
			//now2 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)
			sTime = &now
			orderData.StartTime.Time = now
			orderData.StartTime.Valid = true
			orderData.ExpiryTime.Time = now.Add(time.Duration((orderData.BuyUnit*orderData.ChatUnitMinute)+orderkey.TextOrderExtraAddMinute) * time.Minute)
			orderData.ExpiryTime.Valid = true
			exTime = &(orderData.ExpiryTime.Time)

		case orderkey.ListenerOrderTypeVoiceChat:
			orderData.ExpiryTime.Time = now.Add(orderkey.VoiceChatOrderExpireHour * time.Hour)
			orderData.ExpiryTime.Valid = true
			exTime = &(orderData.ExpiryTime.Time)

		default:
			return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "chat type error")
		}
		// 累计购买该XXX分钟数
		// 查询该对用户历史统计数据
		var buyMinSum int64
		buyMinSum, err = l.getOrderStat(orderData.Uid, orderData.ListenerUid)
		if err != nil {
			return err
		}

		err = l.svcCtx.ChatOrderModel.UpdateOrderPaySuccess(l.ctx, orderData.OrderId, action, buyMinSum, sTime, exTime, states)
		if err != nil {
			return err
		}

		if orderData.ActualAmount > 0 {
			err = l.svcCtx.ChatOrderStatModel.Update2(l.ctx, fmt.Sprintf(db.DBUidId, orderData.Uid, orderData.ListenerUid), orderData.BuyUnit*orderData.ChatUnitMinute)
			if err != nil {
				return err
			}
		}

	default:
		err = l.svcCtx.ChatOrderModel.UpdateOrderState(l.ctx, in.OrderId, action, states)
		if err != nil {
			return err
		}
	}
	return
}

func (l *DoChatOrderActionLogic) getOrderStat(uid, listenerUid int64) (buyMinSum int64, err error) {
	id := fmt.Sprintf(db.DBUidId, uid, listenerUid)

	var data *orderPgModel.ChatOrderStat
	data, err = l.svcCtx.ChatOrderStatModel.FindOne(l.ctx, id)
	if err != nil && err != orderPgModel.ErrNotFound {
		return
	}
	if data == nil && err == orderPgModel.ErrNotFound { // 首次
		data = &orderPgModel.ChatOrderStat{
			Id:               id,
			ListenerUid:      listenerUid,
			Uid:              uid,
			ConfirmMinuteSum: 0,
		}
		_, err = l.svcCtx.ChatOrderStatModel.Insert(l.ctx, data)
		if err != nil {
			return
		}
		buyMinSum = 0
		return
	}

	buyMinSum = data.ConfirmMinuteSum

	return
}

func (l *DoChatOrderActionLogic) verifyState(action int64, orderData *orderPgModel.ChatOrder) (err error, states []int64) {
	switch action {
	case orderkey.ChatOrderStateWaitPay1: // 待支付
		err = errTradeState

	case orderkey.ChatOrderStateCancel2, orderkey.ChatOrderStatePayFail22: // 已取消 订单结束
		states = []int64{orderkey.ChatOrderStateWaitPay1}

	case orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateTextOrderPaySuccess24: // 支付成功 待使用
		states = []int64{orderkey.ChatOrderStateWaitPay1, orderkey.ChatOrderStateCancel2}

	case orderkey.ChatOrderStateUsing4: // 开始使用 服务中
		states = []int64{orderkey.ChatOrderStatePaySuccess3, orderkey.ChatOrderStateUsing4}

	case orderkey.ChatOrderStateExpire17: // 过期
		states = orderkey.CanExpiryOrderSate

	case orderkey.ChatOrderStateApplyRefund5: // 申请退款
		states = orderkey.CanApplyRefundOrderState
		if orderData.ActualAmount <= 0 {
			err = xerr.NewGrpcErrCodeMsg(xerr.OrderError, "订单金额为0，无法申请退款")
			return
		}

	case orderkey.ChatOrderStateSecondApplyRefund10: // 再次申请退款
		states = []int64{orderkey.ChatOrderStateListenerRefuseRefund9}

	case orderkey.ChatOrderStateAutoAgreeRefund6: // 符合自动退款条件
		// 已经判断是否满足自动退款 满足自动同意
		states = []int64{orderkey.ChatOrderStateApplyRefund5}

	case orderkey.ChatOrderStateListenerAgreeRefund7, orderkey.ChatOrderStateListenerRefuseRefund9,
		orderkey.ChatOrderStateAutoAgreeNotProcessRefund27: // XXX同意、拒绝退款、超时未处理系统自动同意退款 退款中
		states = []int64{orderkey.ChatOrderStateApplyRefund5}

	case orderkey.ChatOrderStateFinishRefund8, orderkey.ChatOrderStateRefundPayFail23: // 退款完成/退款失败
		states = orderkey.CanRefundPassOrderState

	case orderkey.ChatOrderStateAdminAgreeRefund11: // 管理员同意退款
		states = []int64{orderkey.ChatOrderStateSecondApplyRefund10}

	case orderkey.ChatOrderStateAutoStartRefund25, orderkey.ChatOrderStateAdminStartRefund26: // 开始退款
		states = orderkey.CanStartRefundOrderState

	case orderkey.ChatOrderStateUserStopService12, orderkey.ChatOrderStateUseOutWaitUserConfirm13: // 用户主动结束 自然结束
		states = orderkey.CanStopOrderState

	case orderkey.ChatOrderStateNot5StarWaitConfirm15: // 用户不满意或者一般评价
		states = orderkey.CanCommentOrderState
		if orderData.CommentTime.Valid {
			err = xerr.NewGrpcErrCodeMsg(xerr.OrderErrorAlreadyComment, "已经评价，无法修改") // 已经评价
			return
		}

	case orderkey.ChatOrderState5StartRatingAndFinish14:
		states = orderkey.CanConfirmFinishOrderState
		if orderData.CommentTime.Valid {
			err = xerr.NewGrpcErrCodeMsg(xerr.OrderErrorAlreadyComment, "已经评价，无法修改") // 已经评价
			return
		}

	case orderkey.ChatOrderStateConfirmFinish16, orderkey.ChatOrderStateAutoCommentFinish18, orderkey.ChatOrderStateAutoConfirmFinish19: // 满意并确认、用户手动确认、系统自动评价并确认、自动确认评价不满意或一般
		states = orderkey.CanConfirmFinishOrderState

	case orderkey.ChatOrderStateRefundRefuseAutoFinish20: // 退款拒绝自动确认
		states = orderkey.RefuseRefundOrderState

	default:
		err = errTradeState
	}
	if !tool.IsInt64ArrayExist(orderData.OrderState, states) {
		err = errTradeState
	}
	return
}

// 是否增加复购人数
func (l *DoChatOrderActionLogic) isRepeatUser(uid, listenerUid int64) (bool, error) {
	// 判断是否是复购用户
	cnt, err := l.svcCtx.ChatOrderModel.FindCount2(l.ctx, uid, listenerUid, 0, false, orderkey.AbnormalOrderState, 0)
	if err != nil {
		return false, err
	}
	if cnt == 1 { // 第二次下单
		return true, nil
	}
	return false, nil
}

// 是否增加服务人数
func (l *DoChatOrderActionLogic) isNewUser(uid, listenerUid int64) (bool, error) {
	// 判断是否增加服务人数
	cnt, err := l.svcCtx.ChatOrderModel.FindCount2(l.ctx, uid, listenerUid, 0, true, orderkey.AbnormalOrderState, 0)
	if err != nil {
		return false, err
	}
	if cnt == 1 { // 第二次下单
		return true, nil
	}
	return false, nil
}
