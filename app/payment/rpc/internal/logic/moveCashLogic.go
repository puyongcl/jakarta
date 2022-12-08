package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/money"
	hfbfcash2 "jakarta/common/third_party/hfbfcash"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
)

type MoveCashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMoveCashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveCashLogic {
	return &MoveCashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  银行卡转账
func (l *MoveCashLogic) MoveCash(in *pb.MoveCashReq) (*pb.MoveCashResp, error) {
	// 查询记录 是否已经提交
	var data *paymentPgModel.ThirdCashFlow
	var err error
	data, err = l.svcCtx.ThirdCashFlowModel.FindOneByWalletFlowNo(l.ctx, in.FlowNo)
	if err != nil && err != paymentPgModel.ErrNotFound {
		return nil, err
	}

	if data != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentErrorFlowAlreadyExist, "提现流水已经存在")
	}
	err = nil

	// 创建记录
	data = &paymentPgModel.ThirdCashFlow{
		FlowNo:       uniqueid.GenSn(uniqueid.SnPrefixThirdCashFlowNo),
		WalletFlowNo: in.FlowNo,
		Amount:       in.Amount,
		PhoneNumber:  in.PhoneNumber,
		Uid:          in.Uid,
		Name:         in.Name,
		IdNo:         in.IdNo,
		BankCardNo:   in.BankCardNo,
		PayStatus:    hfbfcash2.CashStatusCreated,
	}

	// 请求第三方转账
	amount := money.GetYuan(in.Amount)
	req := hfbfcash2.QuickWorkReq{
		Industry:   hfbfcash2.Industry,
		WorkName:   in.FlowNo,
		TotalPrice: amount,
		Detail: []hfbfcash2.QuickWorkDetail{hfbfcash2.QuickWorkDetail{
			Name:         in.Name,
			Idcard:       in.IdNo,
			Bank:         "",
			BankCard:     in.BankCardNo,
			Phone:        in.PhoneNumber,
			Price:        amount,
			CustomNumber: data.FlowNo,
		},
		},
	}

	if l.svcCtx.Config.Mode != service.ProMode {
		return &pb.MoveCashResp{}, nil
	}

	var rsp *hfbfcash2.QuickWorkRsp
	rsp, err = l.svcCtx.HfbfCashClient.QuickWork(&req)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("MoveCashLogic MoveCash QuickWork req:%+v, err:%+v", req, err)
		err = nil
		data.ErrMsg = fmt.Sprintf("%+v", err)
	}
	data.WorkNumber = rsp.Data.WorkNumber

	_, err = l.svcCtx.ThirdCashFlowModel.Insert(l.ctx, data)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("MoveCashLogic MoveCash Insert data:%+v, err:%+v", data, err)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("%+v", err))
	}

	resp := &pb.MoveCashResp{
		CashFlowNo: data.FlowNo,
		Code:       rsp.Code,
		Msg:        rsp.Message,
	}
	return resp, nil
}
