package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerCashLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerCashLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerCashLogLogic {
	return &GetListenerCashLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取提现记录
func (l *GetListenerCashLogLogic) GetListenerCashLog(in *pb.GetListenerCashLogReq) (*pb.GetListenerCashLogResp, error) {
	var settleType []int64
	if in.SettleType == 0 {
		settleType = listenerkey.CashSettleType
	} else {
		settleType = append(settleType, in.SettleType)
	}

	rs, err := l.svcCtx.ListenerWalletFlowModel.FindMoveCashList(l.ctx, in.ListenerUid, settleType, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetListenerCashLogResp{List: make([]*pb.ListenerCashLog, 0)}
	if len(rs) <= 0 {
		return resp, nil
	}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.ListenerCashLog
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].CreateTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
