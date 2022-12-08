package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/payment/mq/internal/svc"
	pbPayment "jakarta/app/payment/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	hfbfcash2 "jakarta/common/third_party/hfbfcash"
	"jakarta/common/third_party/tim"
	"jakarta/common/xerr"
	"time"
)

type UpdateCashStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCashStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCashStatusMq {
	return &UpdateCashStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCashStatusMq) Consume(_, val string) error {
	var message kqueue.UpdateCashStatusMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCashStatusMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCashStatusMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}
	return nil
}

func (l *UpdateCashStatusMq) execService(message *kqueue.UpdateCashStatusMessage) error {
	logx.WithContext(l.ctx).Infof("UpdateCashStatusMq message:%+v", message)

	var payStatus int64
	var customNumber string
	if message.Type == hfbfcash2.CallbackTypePay {
		payStatus = hfbfcash2.GetCashPayStatus(message.PayStatus)
		customNumber = message.CustomNumber
	} else {
		payStatus = hfbfcash2.CashStatusDrop
		if len(message.CustomerNumbers) >= 1 {
			customNumber = message.CustomerNumbers[0]
		}
	}
	payTime, err := time.Parse(hfbfcash2.DateTimeFormat, message.PayTime)
	if err != nil {
		return err
	}

	req := pbPayment.UpdateMoveCashStatusReq{
		WorkNumber:        message.WorkNumber,
		CompanyId:         message.CompanyId,
		TransactionNumber: message.Number,
		PayStatus:         payStatus,
		FlowNo:            customNumber,
		ErrMsg:            message.Msg,
		PayTime:           payTime.Format(db.DateTimeFormat),
	}
	rsp, err := l.svcCtx.PaymentRpc.UpdateMoveCashStatus(l.ctx, &req)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCashStatusMq UpdateMoveCashStatus err:%+v req:%+v", err, message)
		return nil
	}

	// 最终状态 更新钱包流水
	switch payStatus {
	case hfbfcash2.CashStatusSettleSuccess, hfbfcash2.CashStatusSettleFail, hfbfcash2.CashStatusDrop:
		return l.updateWallet(rsp.Uid, rsp.WalletFlowNo, payStatus, customNumber, message.Msg, payTime.Format(db.DateTimeFormat))
	default:
		return nil
	}
}

func (l *UpdateCashStatusMq) updateWallet(uid int64, walletFlowNo string, payStatus int64, cashFlowNo string, msg string, payTime string) error {
	var settleType int64
	switch payStatus {
	case hfbfcash2.CashStatusSettleSuccess:
		settleType = listenerkey.ListenerSettleTypeCashSuccess
	case hfbfcash2.CashStatusSettleFail:
		settleType = listenerkey.ListenerSettleTypeCashFail
	case hfbfcash2.CashStatusDrop:
		settleType = listenerkey.ListenerSettleTypeCashCancel
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("wrong status %d", payStatus))
	}
	var in pbListener.UpdateListenerWalletReq
	in = pbListener.UpdateListenerWalletReq{
		ListenerUid: uid,
		OutId:       cashFlowNo,
		SettleType:  settleType,
		Remark:      msg,
		OutTime:     payTime,
		FlowNo:      walletFlowNo,
	}
	_, err := l.svcCtx.ListenerRpc.UpdateListenerWallet(l.ctx, &in)
	if err != nil {
		return err
	}

	// 消息通知
	l.notify(settleType, uid, msg)
	return nil
}

func (l *UpdateCashStatusMq) notify(status int64, uid int64, errMsg string) {
	var kqmsg *kqueue.SendImDefineMessage
	switch status {
	case listenerkey.ListenerSettleTypeCashSuccess:
		kqmsg = &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   uid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg24,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle24,
			Text:    notify.DefineNotifyMsgTemplateSystemMsg24,
			Val1:    "",
			Val2:    "",
			Sync:    tim.TimMsgSyncFromNo,
		}
	case listenerkey.ListenerSettleTypeCashCancel, listenerkey.ListenerSettleTypeCashFail:
		kqmsg = &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   uid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg25,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle25,
			Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsg25, errMsg),
			Val1:    "",
			Val2:    "",
			Sync:    tim.TimMsgSyncFromNo,
		}
	default:
		return
	}

	buf, err := json.Marshal(kqmsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCashStatusMq json marshal err:%+v", err)
		return
	}

	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCashStatusMq kafka Push err:%+v", err)
		return
	}
	return
}
