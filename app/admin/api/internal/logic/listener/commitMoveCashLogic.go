package listener

import (
	"context"
	"fmt"
	"jakarta/app/admin/api/internal/logic/adminlog"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	pbListener "jakarta/app/listener/rpc/pb"
	pbPayment "jakarta/app/payment/rpc/pb"
	"jakarta/common/key/listenerkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommitMoveCashLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommitMoveCashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommitMoveCashLogic {
	return &CommitMoveCashLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommitMoveCashLogic) CommitMoveCash(req *types.CommitMoveCashReq) (resp *types.CommitMoveCashResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "CheckRefundOrder", req.AdminUid, err, req, resp)
	}()

	for k, _ := range req.Data {
		err = l.commit(req.Data[k])
		if err != nil {
			return
		}
	}
	return
}

func (l *CommitMoveCashLogic) commit(req *types.MoveCashData) error {
	// 获取用户银行卡并更新钱包 TODO 如果报错需要先查询错误并手动处理 确认具体原因 此处不要去处理错误 更新钱包 以免造成金额加错
	var in pbListener.GetCommitMoveCashReq
	in = pbListener.GetCommitMoveCashReq{
		FlowNo: req.FlowNo,
		Uid:    req.Uid,
	}
	rsp1, err := l.svcCtx.ListenerRpc.GetCommitMoveCash(l.ctx, &in)
	if err != nil {
		return err
	}

	// 提交转账
	var in2 pbPayment.MoveCashReq
	in2 = pbPayment.MoveCashReq{
		FlowNo:      req.FlowNo,
		Amount:      rsp1.Amount,
		PhoneNumber: rsp1.PhoneNumber,
		Name:        rsp1.Name,
		IdNo:        rsp1.IdNo,
		BankCardNo:  rsp1.BankCardNo,
		Uid:         rsp1.Uid,
	}

	var in3 pbListener.UpdateListenerWalletReq
	in3 = pbListener.UpdateListenerWalletReq{
		SettleType: listenerkey.ListenerSettleTypeStartCash,
		FlowNo:     in.FlowNo,
	}

	rsp2, err := l.svcCtx.PaymentRpc.MoveCash(l.ctx, &in2)
	if err != nil {
		in3.SettleType = listenerkey.ListenerSettleTypeCashFail
		in3.Remark = fmt.Sprintf("%+v", err)
	} else {
		in3.OutId = rsp2.CashFlowNo
		in3.Remark = fmt.Sprintf("%d:%s", rsp2.Code, rsp2.Msg)
	}

	// 更新钱包流水
	_, err = l.svcCtx.ListenerRpc.UpdateListenerWallet(l.ctx, &in3)
	if err != nil {
		return err
	}
	return nil
}
