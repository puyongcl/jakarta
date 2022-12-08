package logic

import (
	"context"
	"fmt"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/money"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerWalletLogic {
	return &UpdateListenerWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  XXX钱包金额更新
func (l *UpdateListenerWalletLogic) UpdateListenerWallet(in *pb.UpdateListenerWalletReq) (*pb.UpdateListenerWalletResp, error) {
	// 校验
	switch in.SettleType {
	case listenerkey.ListenerSettleTypeOrderAmount:
		if in.OutTime == "" {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "out time empty")
		}

	case listenerkey.ListenerSettleTypeApplyCash: // 申请提现
		// 校验是否是提现时间段
		day := time.Now().Day()
		if day < listenerkey.ListenerMoveCashStatDay || day > listenerkey.ListenerMoveCashEndDay {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorMoveCash, fmt.Sprintf("请在每月%d-%d号申请提现", listenerkey.ListenerMoveCashStatDay, listenerkey.ListenerMoveCashEndDay))
		}

		// 校验银行卡是否绑定
		_, err := l.svcCtx.ListenerBankCardModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil && err != listenerPgModel2.ErrNotFound {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, err.Error())
		}
		if err == listenerPgModel2.ErrNotFound {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorNotSetBankCard, "请设置结算方式")
		}

		// 获取当前钱包余额
		if in.Amount == 0 {
			var data *listenerPgModel2.ListenerWallet
			data, err = l.svcCtx.ListenerWalletModel.FindOne2(l.ctx, in.ListenerUid)
			if err != nil {
				return nil, err
			}
			in.Amount = data.Amount - data.CurrentMonthAmount

			// 校验提现金额
			if in.Amount < listenerkey.ListenerMoveCashMinAmount {
				return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorMoveCash, fmt.Sprintf("可提现金额门槛%s元，请先接单获得XXX收益", money.GetYuan(listenerkey.ListenerMoveCashMinAmount)))
			}
		}

	default:

	}

	// 更新

	var err error
	switch in.SettleType {
	case listenerkey.ListenerSettleTypeConfirm: // 已确认收益
		var outTime time.Time
		outTime, err = time.ParseInLocation(db.DateTimeFormat, in.OutTime, time.Local)
		if err != nil {
			return nil, err
		}
		var addMount int64 // 统计本月订单的金额
		if outTime.Month() == time.Now().Month() {
			addMount = in.Amount
		}
		_, err = l.svcCtx.ListenerWalletFlowModel.UpdateConfirmIncomeLog(l.ctx, in.OutId, in.SettleType, in.Amount, in.Remark)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.AddConfirmAmount(l.ctx, in.ListenerUid, in.Amount, addMount)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeOrderAmount: // 本月订单总金额
		var outTime time.Time
		outTime, err = time.ParseInLocation(db.DateTimeFormat, in.OutTime, time.Local)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.AddOrderAmount(l.ctx, in.ListenerUid, in.OrderAmount)
		if err != nil {
			return nil, err
		}
		_, err = l.svcCtx.ListenerWalletFlowModel.InsertIncomeLog(l.ctx, &listenerPgModel2.ListenerWalletFlow{
			FlowNo:      uniqueid.GenSn(uniqueid.SnPrefixWalletFlowNo),
			ListenerUid: in.ListenerUid,
			OrderAmount: in.OrderAmount,
			OutId:       in.OutId,
			SettleType:  in.SettleType,
			Remark:      in.Remark,
			OutTime:     outTime,
		})
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeAlreadyRefund: // 已经退款
		_, err = l.svcCtx.ListenerWalletFlowModel.UpdateRefundIncomeLog(l.ctx, in.OutId, in.SettleType, in.Remark)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.RefundAmount(l.ctx, in.ListenerUid, in.OrderAmount)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeApplyCash: // 申请提现
		_, err = l.svcCtx.ListenerWalletFlowModel.InsertOutLogTrans(l.ctx, &listenerPgModel2.ListenerWalletFlow{
			FlowNo:      uniqueid.GenSn(uniqueid.SnPrefixWalletFlowNo),
			ListenerUid: in.ListenerUid,
			Amount:      in.Amount,
			OutId:       in.OutId,
			SettleType:  in.SettleType,
			Remark:      in.Remark,
			OutTime:     time.Now(),
		})
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.ApplyCashAmount(l.ctx, in.ListenerUid, in.Amount)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeStartCash: // 开始提现
		_, err = l.svcCtx.ListenerWalletFlowModel.UpdateOutLog(l.ctx, in.FlowNo, in.OutId, in.SettleType, in.Remark)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeCashSuccess: // 提现成功
		_, err = l.svcCtx.ListenerWalletFlowModel.UpdateOutLog(l.ctx, in.FlowNo, in.OutId, in.SettleType, in.Remark)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.AlreadyCashAmount(l.ctx, in.ListenerUid, in.Amount)
		if err != nil {
			return nil, err
		}

	case listenerkey.ListenerSettleTypeCashFail, listenerkey.ListenerSettleTypeCashCancel: // 提现失败
		_, err = l.svcCtx.ListenerWalletFlowModel.UpdateOutLog(l.ctx, in.FlowNo, in.OutId, in.SettleType, in.Remark)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ListenerWalletModel.CashFail(l.ctx, in.ListenerUid, in.Amount)
		if err != nil {
			return nil, err
		}

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "wrong settle type")
	}

	return &pb.UpdateListenerWalletResp{}, nil
}
