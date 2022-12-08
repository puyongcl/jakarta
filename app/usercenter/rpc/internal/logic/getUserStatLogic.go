package logic

import (
	"context"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatLogic {
	return &GetUserStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserStatLogic) GetUserStat(in *pb.GetUserStatReq) (*pb.GetUserStatResp, error) {
	rs, err := l.svcCtx.UserStatModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetUserStatResp{
		Uid:             rs.Uid,
		CostAmountSum:   rs.CostAmountSum,
		RefundAmountSum: rs.RefundAmountSum,
		PaidOrderCnt:    rs.PaidOrderCnt,
		RefundOrderCnt:  rs.RefundOrderCnt,
	}
	return resp, nil
}
