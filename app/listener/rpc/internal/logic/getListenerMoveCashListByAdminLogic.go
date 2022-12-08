package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerMoveCashListByAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerMoveCashListByAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerMoveCashListByAdminLogic {
	return &GetListenerMoveCashListByAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理后台获取XXX提现申请列表
func (l *GetListenerMoveCashListByAdminLogic) GetListenerMoveCashListByAdmin(in *pb.GetListenerMoveCashListByAdminReq) (*pb.GetListenerMoveCashListByAdminResp, error) {
	var settleType []int64
	if in.SettleType == 0 {
		settleType = listenerkey.CashSettleType
	} else {
		settleType = append(settleType, in.SettleType)
	}

	var err error
	resp := &pb.GetListenerMoveCashListByAdminResp{List: make([]*pb.AdminSeeListenerMoveCash, 0)}
	resp.Sum, err = l.svcCtx.ListenerWalletFlowModel.CountMoveCashList(l.ctx, in.ListenerUid, settleType)
	if err != nil {
		return nil, err
	}
	//
	rs, err := l.svcCtx.ListenerWalletFlowModel.FindMoveCashList(l.ctx, in.ListenerUid, settleType, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rs) <= 0 {
		return resp, nil
	}

	for idx := 0; idx < len(rs); idx++ {
		var val pb.AdminSeeListenerMoveCash
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].CreateTime.Format(db.DateTimeFormat)

		var lbc *listenerPgModel.ListenerBankCard
		lbc, err = l.svcCtx.ListenerBankCardModel.FindOne(l.ctx, rs[idx].ListenerUid)
		if err != nil {
			return nil, err
		}
		val.ListenerName = lbc.ListenerName
		val.IdNo = lbc.IdNo
		val.BankCardNo = lbc.BankCardNo

		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
