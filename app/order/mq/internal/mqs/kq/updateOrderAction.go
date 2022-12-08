package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/order/mq/internal/svc"
	pbOrder "jakarta/app/order/rpc/pb"
	pbPayment "jakarta/app/payment/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"
	"jakarta/common/kqueue"
	"jakarta/common/money"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strconv"
	"time"
)

/**
Listening to the payment flow status change notification message queue
*/
type UpdateOrderActionMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderActionMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderActionMq {
	return &UpdateOrderActionMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderActionMq) Consume(_, val string) error {
	var message kqueue.UpdateChatOrderActionMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateOrderActionMq->Consume Unmarshal err : %+v , val: %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateOrderActionMq->execService err : %+v , val :%s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *UpdateOrderActionMq) execService(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	logx.WithContext(l.ctx).Infof("UpdateOrderActionMq time:%s msg:%+v", time.Now().Format(db.DateTimeFormat), message)
	switch message.Action {
	case orderkey.ChatOrderStateCancel2: // 已取消 订单结束

	case orderkey.ChatOrderStatePaySuccess3: // 支付成功 待使用
		return l.paySuccess(message)

	case orderkey.ChatOrderStatePayFail22: // 支付失败

	case orderkey.ChatOrderStateRefundPayFail23: // 退款失败

	case orderkey.ChatOrderStateUsing4: // 开始使用 服务中

	case orderkey.ChatOrderStateTextOrderPaySuccess24: // 文字订单支付成功 即时生效开始服务
		return l.paySuccess(message)

	case orderkey.ChatOrderStateExpire17: // 过期
		return l.expiry(message)

	case orderkey.ChatOrderStateApplyRefund5: // 申请退款
		return l.applyRefund(message)

	case orderkey.ChatOrderStateAutoAgreeRefund6: // 自动退款 退款同意之后 需要延后等待两天 定时任务发起真正退款
		return l.autoRefund(message)

	case orderkey.ChatOrderStateAutoAgreeNotProcessRefund27: // 系统自动同意退款 XXX超时未处理
		return l.listenerAgreeRefund(message)

	case orderkey.ChatOrderStateListenerAgreeRefund7: // XXX同意退款
		return l.listenerAgreeRefund(message)

	case orderkey.ChatOrderStateListenerRefuseRefund9: // XXX拒绝退款
		return l.listenerRefuseRefund(message)

	case orderkey.ChatOrderStateAdminAgreeRefund11: // 管理员同意退款
		return l.adminAgreeRefund(message)

	case orderkey.ChatOrderStateUserStopService12: // 用户主动结束
		return l.userStop(message)

	case orderkey.ChatOrderStateUseOutWaitUserConfirm13: // 自然结束
		return l.useOut(message)

	case orderkey.ChatOrderStateNot5StarWaitConfirm15: // 用户不满意或者一般评价
		return l.lowerComment(message)

	case orderkey.ChatOrderStateAutoStartRefund25, orderkey.ChatOrderStateAdminStartRefund26: // 开始退款
		return l.refund(message)

	case orderkey.ChatOrderStateFinishRefund8: // 退款完成
		return l.finishRefund(message)

	case orderkey.ChatOrderState5StartRatingAndFinish14: // 好评自动确认完成
		return l.goodComment(message)

	case orderkey.ChatOrderStateAutoCommentFinish18, orderkey.ChatOrderStateRefundRefuseAutoFinish20,
		orderkey.ChatOrderStateAutoConfirmFinish19, orderkey.ChatOrderStateConfirmFinish16: // 确认完成、退款拒绝 进行结算
		return l.finishSettle(message)

	default:
		return nil
	}
	return nil
}

func (l *UpdateOrderActionMq) userStop(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	if message.UsedMinute < message.BuyMinute {
		in := pbChat.UpdateUserChatBalanceReq{
			Uid:         message.Uid,
			ListenerUid: message.ListenerUid,
			OrderType:   message.OrderType,
			AddMinute:   message.BuyMinute - message.UsedMinute,
			OrderId:     message.OrderId,
			Action:      message.Action,
			EventType:   chatkey.ChatStatUpdateTypeOrderUserStop,
		}
		if message.OrderType == orderkey.ListenerOrderTypeTextChat {
			in.AddMinute += orderkey.TextOrderExtraAddMinute
		}

		err = l.updateChatBalance(&in)
		if err != nil {
			return
		}
	}

	if message.UsedMinute > 0 {
		_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &pbListener.UpdateListenerOrderStatReq{
			AddUserCount:    message.AddUser,
			AddChatDuration: message.UsedMinute,
			ListenerUid:     message.ListenerUid,
		})
		if err != nil {
			return
		}
	}

	l.resetFreeChatCnt(message.Uid, message.ListenerUid)
	return
}

func (l *UpdateOrderActionMq) resetFreeChatCnt(uid, listenerUid int64) {
	_, err := l.svcCtx.ChatRpc.ResetFreeTextChatCnt(l.ctx, &pbChat.ResetFreeTextChatCntReq{
		Uid:         uid,
		ListenerUid: listenerUid,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateOrderActionMq ResetFreeTextChatCnt err:%+v", err)
		return
	}
}

func (l *UpdateOrderActionMq) useOut(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	l.resetFreeChatCnt(message.Uid, message.ListenerUid)

	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &pbListener.UpdateListenerOrderStatReq{
		AddUserCount:    message.AddUser,
		AddChatDuration: message.UsedMinute,
		ListenerUid:     message.ListenerUid,
	})
	if err != nil {
		return
	}
	return
}

func (l *UpdateOrderActionMq) finishSettle(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	var amount int64
	amount, err = l.settleOrder(message, listenerkey.ListenerSettleTypeConfirm)
	if err != nil {
		return
	}
	// 消息通知
	if message.Action == orderkey.ChatOrderStateRefundRefuseAutoFinish20 { // 再次申请退款或者XXX拒绝退款 1天后 自动确认 给用户通知客服拒绝
		err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg12,
			notify.DefineNotifyMsgTemplateOrderRefundMsgTittle12, notify.DefineNotifyMsgTemplateOrderRefundMsg12, message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return
		}
	}
	// 给XXX 收益到账
	if amount > 0 {
		err = l.notify(notify.TimOrderNotifyUid, message.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg16, notify.DefineNotifyMsgTemplateOrderMsgTittle16, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg16, message.OrderId, money.GetYuan(amount)), message.OrderId, strconv.FormatInt(message.Uid, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return
		}
	}

	// 更新XXX服务统计
	var in = pbListener.UpdateListenerOrderStatReq{
		AddFinishOrderCnt: 1,
		ListenerUid:       message.ListenerUid,
	}
	if message.Action == orderkey.ChatOrderStateAutoCommentFinish18 { // 自动好评
		in.AddFiveStar = 1
		in.AddRatingSum = 1
	}
	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &in)
	if err != nil {
		return
	}
	return
}

func (l *UpdateOrderActionMq) goodComment(msg *kqueue.UpdateChatOrderActionMessage) (err error) {
	var amount int64
	amount, err = l.settleOrder(msg, listenerkey.ListenerSettleTypeConfirm)
	if err != nil {
		return
	}

	err = l.notify(notify.TimOrderNotifyUid, msg.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg5,
		notify.DefineNotifyMsgTemplateOrderMsgTitle5, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg5,
			msg.NickName, listenerkey.GetRatingText(msg.Star)), msg.OrderId,
		strconv.FormatInt(msg.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}
	// 同步评价到聊天
	if msg.SendMsg == db.Enable && msg.Comment != "" {
		// 发送给XXX
		err = l.notify(msg.Uid, msg.ListenerUid, notify.DefineNotifyMsgTypeChatMsg24,
			"", msg.Comment, strconv.FormatInt(msg.Star, 10),
			tool.GetIntArrayString(msg.CommentTag), tim.TimMsgSyncFromYes)
		if err != nil {
			return
		}
	}
	// 发送给用户 感谢评价
	err = l.notify(msg.ListenerUid, msg.Uid, notify.DefineNotifyMsgTypeChatMsg26,
		"", notify.DefineNotifyMsgTemplateChatMsg26, "",
		"", tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}

	// 给XXX 收益到账
	err = l.notify(notify.TimOrderNotifyUid, msg.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg16, notify.DefineNotifyMsgTemplateOrderMsgTittle16, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg16, msg.OrderId, money.GetYuan(amount)), msg.OrderId, strconv.FormatInt(msg.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}
	// 更新XXX服务情况
	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &pbListener.UpdateListenerOrderStatReq{
		AddRatingSum:      1,
		AddFiveStar:       1,
		AddFinishOrderCnt: 1,
		ListenerUid:       msg.ListenerUid,
		CommentTag:        msg.CommentTag,
	})
	if err != nil {
		return
	}
	// 发送用户和XXX互动事件
	et := pbChat.SendUserListenerRelationEventReq{
		Uid:         msg.Uid,
		ListenerUid: msg.ListenerUid,
		EventType:   chatkey.InteractiveEventTypeCommentOrder5Star5,
	}
	_, _ = l.svcCtx.ChatRpc.SendUserListenerRelationEvent(l.ctx, &et)
	return
}

func (l *UpdateOrderActionMq) finishRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	_, err = l.settleOrder(message, listenerkey.ListenerSettleTypeAlreadyRefund)
	if err != nil {
		return
	}
	err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg15,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle15, notify.DefineNotifyMsgTemplateOrderRefundMsg15,
		message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}
	// 更新XXX统计数据
	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &pbListener.UpdateListenerOrderStatReq{
		AddRefundOrderCnt: 1,
		ListenerUid:       message.ListenerUid,
	})
	if err != nil {
		return
	}
	// 更新用户统计数据
	_, err = l.svcCtx.UsercenterRpc.UpdateUserStat(l.ctx, &pbUser.UpdateUserStatReq{
		Uid:                message.Uid,
		AddRefundAmountSum: message.PaidAmount,
		AddRefundOrderCnt:  1,
	})
	return
}

func (l *UpdateOrderActionMq) lowerComment(msg *kqueue.UpdateChatOrderActionMessage) (err error) {
	err = l.notify(notify.TimOrderNotifyUid, msg.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg5,
		notify.DefineNotifyMsgTemplateOrderMsgTitle5, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg5,
			msg.NickName, listenerkey.GetRatingText(msg.Star)), msg.OrderId,
		strconv.FormatInt(msg.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}
	// 同步评价到聊天
	if msg.SendMsg == db.Enable && msg.Comment != "" {
		err = l.notify(msg.Uid, msg.ListenerUid, notify.DefineNotifyMsgTypeChatMsg24,
			"", msg.Comment, strconv.FormatInt(msg.Star, 10),
			tool.GetIntArrayString(msg.CommentTag), tim.TimMsgSyncFromYes)
		if err != nil {
			return
		}
	}
	// 发送给用户 感谢评价
	err = l.notify(msg.ListenerUid, msg.Uid, notify.DefineNotifyMsgTypeChatMsg26,
		"", notify.DefineNotifyMsgTemplateChatMsg26, "",
		"", tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}

	in := pbListener.UpdateListenerOrderStatReq{
		AddRatingSum: 1,
		ListenerUid:  msg.ListenerUid,
		CommentTag:   msg.CommentTag,
	}
	if msg.Star == listenerkey.Rating3Star {
		in.AddThreeStar = 1
	} else if msg.Star == listenerkey.Rating1Star {
		in.AddOneStar = 1
	}
	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &in)
	if err != nil {
		return
	}

	// 发送用户和XXX互动事件
	et := pbChat.SendUserListenerRelationEventReq{
		Uid:         msg.Uid,
		ListenerUid: msg.ListenerUid,
	}
	if msg.Star == listenerkey.Rating3Star {
		et.EventType = chatkey.InteractiveEventTypeCommentOrder3Star6
	} else if msg.Star == listenerkey.Rating1Star {
		et.EventType = chatkey.InteractiveEventTypeCommentOrder1Star7
	}
	_, _ = l.svcCtx.ChatRpc.SendUserListenerRelationEvent(l.ctx, &et)
	return
}

func (l *UpdateOrderActionMq) listenerRefuseRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg10,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle10, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderRefundMsg10, message.Reason),
		message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return
	}
	return
}

func (l *UpdateOrderActionMq) adminAgreeRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	// 给用户通知
	err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg13,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle13, notify.DefineNotifyMsgTemplateOrderRefundMsg13,
		message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	// 给XXX
	orderCreateTime, err := time.ParseInLocation(db.DateTimeFormat, message.OrderCreateTime, time.Local)
	if err != nil {
		return err
	}
	err = l.notify(notify.TimOrderNotifyUid, message.ListenerUid, notify.DefineNotifyMsgTypeOrderRefundMsg14,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle14, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderRefundMsg14,
			message.NickName, orderCreateTime.Month(), orderCreateTime.Day(), message.OrderId, message.Reason),
		message.OrderId, strconv.FormatInt(message.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	return
}

func (l *UpdateOrderActionMq) listenerAgreeRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg11,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle11, notify.DefineNotifyMsgTemplateOrderRefundMsg11, message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	return
}

func (l *UpdateOrderActionMq) autoRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	// 对用户
	err = l.notify(notify.TimSystemNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderRefundMsg7,
		notify.DefineNotifyMsgTemplateOrderMsgTitle7, notify.DefineNotifyMsgTemplateOrderRefundMsg7,
		message.OrderId, strconv.FormatInt(message.ListenerUid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	// 给XXX
	orderCreateTime, err := time.ParseInLocation(db.DateTimeFormat, message.OrderCreateTime, time.Local)
	if err != nil {
		return err
	}
	err = l.notify(notify.TimSystemNotifyUid, message.ListenerUid, notify.DefineNotifyMsgTypeOrderRefundMsg8,
		notify.DefineNotifyMsgTemplateOrderMsgTitle8, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderRefundMsg8,
			message.NickName, orderCreateTime.Month(), orderCreateTime.Day(), message.OrderId, message.Reason),
		message.OrderId, strconv.FormatInt(message.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	return
}

func (l *UpdateOrderActionMq) applyRefund(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	// 给XXX
	orderCreateTime, err := time.ParseInLocation(db.DateTimeFormat, message.OrderCreateTime, time.Local)
	if err != nil {
		return err
	}
	err = l.notify(notify.TimOrderNotifyUid, message.ListenerUid, notify.DefineNotifyMsgTypeOrderRefundMsg9,
		notify.DefineNotifyMsgTemplateOrderRefundMsgTittle9, fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderRefundMsg9,
			message.NickName, orderCreateTime.Month(), orderCreateTime.Day(), message.OrderId, message.Reason),
		message.OrderId, strconv.FormatInt(message.Uid, 10), tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	return
}
func (l *UpdateOrderActionMq) expiry(message *kqueue.UpdateChatOrderActionMessage) (err error) {
	// 更新可用时间
	if message.BuyMinute > message.UsedMinute {
		err = l.updateChatBalance(&pbChat.UpdateUserChatBalanceReq{
			Uid:         message.Uid,
			ListenerUid: message.ListenerUid,
			OrderType:   message.OrderType,
			AddMinute:   -(message.BuyMinute - message.UsedMinute),
			OrderId:     message.OrderId,
			Action:      message.Action,
			EventType:   chatkey.ChatStatUpdateTypeOrderExpireDecr,
		})
		if err != nil {
			return err
		}
	}

	// 发送消息
	err = l.notify(notify.TimOrderNotifyUid, message.Uid, notify.DefineNotifyMsgTypeOrderMsg4,
		notify.DefineNotifyMsgTemplateOrderMsgTitle4, notify.DefineNotifyMsgTemplateOrderMsg4,
		strconv.FormatInt(message.ListenerUid, 10), message.OrderId, tim.TimMsgSyncFromNo)
	if err != nil {
		return err
	}
	return
}

// 增加或减少时间
func (l *UpdateOrderActionMq) updateChatBalance(in *pbChat.UpdateUserChatBalanceReq) error {
	// 更新用户可用时间 chat rpc
	_, err := l.svcCtx.ChatRpc.UpdateUserChatBalance(l.ctx, in)
	if err != nil {
		return xerr.NewGrpcErrCodeMsg(xerr.ServerCommonError, fmt.Sprintf("req:%+v, err:%+v", in, err))
	}
	return nil
}

// 结算
func (l *UpdateOrderActionMq) settleOrder(message *kqueue.UpdateChatOrderActionMessage, settleType int64) (int64, error) {
	var amount int64 // XXX收益
	//var orderAmount int64 // 退款金额 订单金额
	//orderAmount = message.PaidAmount
	switch settleType {
	case listenerkey.ListenerSettleTypeOrderAmount:
		// 待确认 未订单总金额
		amount = message.PaidAmount
	case listenerkey.ListenerSettleTypeConfirm:
		// 更新订单分成
		rs, err := l.svcCtx.OrderRpc.SettleChatOrder(l.ctx, &pbOrder.SettleChatOrderReq{
			OrderId:     message.OrderId,
			Star:        message.Star,
			SettleType:  settleType,
			Uid:         message.Uid,
			ListenerUid: message.ListenerUid,
		})
		if err != nil {
			return 0, err
		}
		amount = rs.Amount

	case listenerkey.ListenerSettleTypeAlreadyRefund:
		// 退款金额为订单总金额
		amount = message.PaidAmount

	default:
		return 0, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("错误的类型:%d", settleType))
	}

	if amount <= 0 {
		return amount, nil
	}

	// 更新XXX钱包
	_, err := l.svcCtx.ListenerRpc.UpdateListenerWallet(l.ctx, &pbListener.UpdateListenerWalletReq{
		ListenerUid: message.ListenerUid,
		Amount:      amount,
		OrderAmount: message.PaidAmount,
		OutId:       message.OrderId,
		SettleType:  settleType,
		OutTime:     message.OrderCreateTime,
	})
	if err != nil {
		return 0, err
	}
	return amount, nil
}

// 退款
func (l *UpdateOrderActionMq) refund(message *kqueue.UpdateChatOrderActionMessage) error {
	_, err := l.svcCtx.PaymentRpc.RequestRefund(l.ctx, &pbPayment.RequestRefundReq{
		OrderId: message.OrderId,
		Reason:  "用户申请退款",
	})
	if err != nil {
		return err
	}
	return nil
}

// 支付成功
func (l *UpdateOrderActionMq) paySuccess(msg *kqueue.UpdateChatOrderActionMessage) (err error) {
	l.resetFreeChatCnt(msg.Uid, msg.ListenerUid)

	// 增加时间
	err = l.updateChatBalance(&pbChat.UpdateUserChatBalanceReq{
		Uid:            msg.Uid,
		ListenerUid:    msg.ListenerUid,
		OrderType:      msg.OrderType,
		AddMinute:      msg.BuyMinute,
		OrderId:        msg.OrderId,
		Action:         msg.Action,
		EventType:      chatkey.ChatStatUpdateTypeOrderPaidAdd,
		TextExpireTime: msg.TextExpireTime,
	})
	if err != nil {
		return err
	}
	// 发消息
	if msg.OrderType == orderkey.ListenerOrderTypeTextChat { // 文字
		err = l.notify2(notify.TimOrderNotifyUid, msg.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg1, notify.DefineNotifyMsgTemplateOrderMsgTitle1,
			fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg1, msg.NickName, msg.BuyMinute), strconv.FormatInt(msg.Uid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo, msg)
		if err != nil {
			return err
		}
		err = l.notify(msg.ListenerUid, msg.Uid, notify.DefineNotifyMsgTypeChatMsg1, "",
			notify.DefineNotifyMsgTemplateChatMsg1, strconv.FormatInt(msg.ListenerUid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return err
		}
		err = l.notify(msg.Uid, msg.ListenerUid, notify.DefineNotifyMsgTypeChatMsg3, "",
			fmt.Sprintf(notify.DefineNotifyMsgTemplateChatMsg3, msg.NickName, msg.BuyMinute), strconv.FormatInt(msg.Uid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return err
		}
	} else if msg.OrderType == orderkey.ListenerOrderTypeVoiceChat { // 语音
		err = l.notify2(notify.TimOrderNotifyUid, msg.ListenerUid, notify.DefineNotifyMsgTypeOrderMsg2, notify.DefineNotifyMsgTemplateOrderMsgTitle2,
			fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg2, msg.NickName, msg.BuyMinute), strconv.FormatInt(msg.Uid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo, msg)
		if err != nil {
			return err
		}
		err = l.notify(msg.ListenerUid, msg.Uid, notify.DefineNotifyMsgTypeChatMsg2, "",
			notify.DefineNotifyMsgTemplateChatMsg2, strconv.FormatInt(msg.ListenerUid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return err
		}
		err = l.notify(msg.Uid, msg.ListenerUid, notify.DefineNotifyMsgTypeChatMsg4, "",
			fmt.Sprintf(notify.DefineNotifyMsgTemplateChatMsg4, msg.NickName, msg.BuyMinute), strconv.FormatInt(msg.Uid, 10), strconv.FormatInt(msg.BuyMinute, 10), tim.TimMsgSyncFromNo)
		if err != nil {
			return err
		}
	}
	// 更新XXX支付成功订单数
	_, err = l.svcCtx.ListenerRpc.UpdateListenerOrderStat(l.ctx, &pbListener.UpdateListenerOrderStatReq{
		AddPaidOrderCnt:      1,
		ListenerUid:          msg.ListenerUid,
		AddRepeatPaidUserCnt: msg.AddRepeat,
		UserPaidOrderAmount:  msg.PaidAmount,
	})
	if err != nil {
		return err
	}

	// 更新用户支付成功统计
	_, err = l.svcCtx.UsercenterRpc.UpdateUserStat(l.ctx, &pbUser.UpdateUserStatReq{
		Uid:              msg.Uid,
		AddCostAmountSum: msg.PaidAmount,
		AddPaidOrderCnt:  1,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateOrderActionMq UpdateUserStat err:%+v", err)
	}
	// 更新待确认收益
	_, err = l.settleOrder(msg, listenerkey.ListenerSettleTypeOrderAmount)
	if err != nil {
		return
	}

	// 发送用户和XXX互动事件
	et := pbChat.SendUserListenerRelationEventReq{
		Uid:         msg.Uid,
		ListenerUid: msg.ListenerUid,
	}
	if msg.OrderType == orderkey.ListenerOrderTypeTextChat {
		et.EventType = chatkey.InteractiveEventTypePayOrder3
	} else {
		et.EventType = chatkey.InteractiveEventTypePayOrder4
	}
	_, err = l.svcCtx.ChatRpc.SendUserListenerRelationEvent(l.ctx, &et)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateOrderActionMq SendUserListenerRelationEvent order id:%s err:%+v", msg.OrderId, err)
	}

	// 上报事件
	l.payEvent(msg.Uid, msg.UserChannel, strconv.FormatInt(msg.PaidAmount, 10))

	logx.WithContext(l.ctx).Infof("UpdateOrderActionMq pay success order id:%s", msg.OrderId)
	return
}

// 消息通知 需要同步发送公众号
func (l *UpdateOrderActionMq) notify2(from, to, msgType int64, title, text, val1, val2 string, sync int64, kqMsg *kqueue.UpdateChatOrderActionMessage) error {
	switch msgType {
	case notify.DefineNotifyMsgTypeOrderMsg1, notify.DefineNotifyMsgTypeOrderMsg2:
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	var orderTypeText string
	if kqMsg.OrderType == orderkey.ListenerOrderTypeVoiceChat {
		orderTypeText = fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg2OrderTyeVoice, kqMsg.BuyMinute)
	} else {
		orderTypeText = fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg2OrderTyeText, kqMsg.BuyMinute)
	}
	msg := &kqueue.SendImDefineMessage{
		FromUid: from,
		ToUid:   to,
		MsgType: msgType,
		Title:   title,
		Text:    text,
		Val1:    val1,
		Val2:    val2,
		Sync:    sync,
		Val3:    kqMsg.OrderCreateTime,
		Val4:    orderTypeText,
		Val5:    kqMsg.NickName,
		Val6:    fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg2OrderInfo, money.GetYuan(kqMsg.PaidAmount)),
	}
	var buf []byte
	var err error
	buf, err = json.Marshal(msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		return err
	}
	return nil
}

// 消息通知
func (l *UpdateOrderActionMq) notify(from, to, msgType int64, title, text, val1, val2 string, sync int64) error {
	// 发送im消息
	msg := &kqueue.SendImDefineMessage{
		FromUid: from,
		ToUid:   to,
		MsgType: msgType,
		Title:   title,
		Text:    text,
		Val1:    val1,
		Val2:    val2,
		Sync:    sync,
	}
	var buf []byte
	var err error
	buf, err = json.Marshal(msg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		return err
	}
	return nil
}

// 上报获客渠道用户付费
func (l *UpdateOrderActionMq) payEvent(uid int64, channel string, value string) {
	var kqMsg *kqueue.UploadUserEventMessage
	switch channel {
	case userkey.GetUserChannelZhihu:
		// 获取用户cb地址
		rsp, err := l.svcCtx.UsercenterRpc.GetUserChannelCallback(l.ctx, &pbUser.GetUserChannelCallbackReq{Uid: uid})
		if err != nil {
			return
		}
		if rsp.Cb == "" {
			return
		}
		//
		kqMsg = &kqueue.UploadUserEventMessage{
			Uid:   uid,
			Cb:    rsp.Cb,
			Event: userkey.ZhihuUploadEventPaid,
			Value: value,
			Stamp: fmt.Sprintf("%d", time.Now().Unix()),
		}
	default:
		return
	}

	buf, err := json.Marshal(kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RegisterLogic regEvent json Marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueUploadUserEventClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RegisterLogic regEvent push err:%+v", err)
		return
	}
	return
}
