package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/db"
	"jakarta/common/third_party/hfbfcash"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMoveCashStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMoveCashStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMoveCashStatusLogic {
	return &UpdateMoveCashStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新转账状态
func (l *UpdateMoveCashStatusLogic) UpdateMoveCashStatus(in *pb.UpdateMoveCashStatusReq) (*pb.UpdateMoveCashStatusResp, error) {
	data, err := l.svcCtx.ThirdCashFlowModel.FindOne(l.ctx, in.FlowNo)
	if err != nil {
		return nil, err
	}
	// 校验状态
	switch data.PayStatus {
	case hfbfcash.CashStatusSettleSuccess, hfbfcash.CashStatusSettleFail, hfbfcash.CashStatusDrop:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("cash status error %d req:%d", data.PayStatus, in.PayStatus))
	default:

	}

	// 更新
	newData := new(paymentPgModel.ThirdCashFlow)
	newData.PayStatus = in.PayStatus
	newData.FlowNo = in.FlowNo
	if in.PayTime != "" {
		newData.PayTime.Time, err = time.ParseInLocation(db.DateTimeFormat, in.PayTime, time.Local)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UpdateMoveCashStatusLogic parse time error:%+v", err)
		} else {
			newData.PayTime.Valid = true
		}
	}
	switch in.PayStatus {
	case hfbfcash.CashStatusCreated: // 已请求
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("cash status error %d req:%d", data.PayStatus, in.PayStatus))
	default:
		if data.WorkNumber == "" {
			newData.WorkNumber = in.WorkNumber
		} else if data.WorkNumber != in.WorkNumber {
			logx.WithContext(l.ctx).Errorf("UpdateMoveCashStatusLogic work number error flowNo:%s work no", in.FlowNo, in.WorkNumber)
			in.ErrMsg += fmt.Sprintf("work number error:%s", in.WorkNumber)
		}
		if in.TransactionNumber != "" {
			newData.TransactionNumber = in.TransactionNumber
		}

		newData.ErrMsg = in.ErrMsg
	}
	err = l.svcCtx.ThirdCashFlowModel.UpdatePart(l.ctx, newData)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateMoveCashStatusResp{Uid: data.Uid, WalletFlowNo: data.WalletFlowNo}, nil
}
