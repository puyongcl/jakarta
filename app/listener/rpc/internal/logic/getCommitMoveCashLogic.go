package logic

import (
	"context"
	"fmt"
	"jakarta/common/key/listenerkey"
	"jakarta/common/xerr"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommitMoveCashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommitMoveCashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommitMoveCashLogic {
	return &GetCommitMoveCashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取提交转账信息并更新状态
func (l *GetCommitMoveCashLogic) GetCommitMoveCash(in *pb.GetCommitMoveCashReq) (*pb.GetCommitMoveCashResp, error) {
	rs1, err := l.svcCtx.ListenerBankCardModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	rs2, err := l.svcCtx.ListenerWalletFlowModel.FindOne(l.ctx, in.FlowNo)
	if err != nil {
		return nil, err
	}
	if rs2.ListenerUid != rs1.ListenerUid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "uid not match")
	}
	if rs2.SettleType != listenerkey.ListenerSettleTypeApplyCash {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorStartMoveCashError, fmt.Sprintf("流水状态错误:%d", rs2.SettleType))
	}

	resp := pb.GetCommitMoveCashResp{
		FlowNo:      in.FlowNo,
		Amount:      rs2.Amount,
		PhoneNumber: rs1.PhoneNumber,
		Name:        rs1.ListenerName,
		IdNo:        rs1.IdNo,
		BankCardNo:  rs1.BankCardNo,
		Uid:         in.Uid,
	}
	return &resp, nil
}
