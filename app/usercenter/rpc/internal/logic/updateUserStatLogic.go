package logic

import (
	"context"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatLogic {
	return &UpdateUserStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserStatLogic) UpdateUserStat(in *pb.UpdateUserStatReq) (*pb.UpdateUserStatResp, error) {
	err := l.svcCtx.UserStatModel.UpdateUserStat(l.ctx, &userPgModel.UpdateUserStatData{
		Uid:                in.Uid,
		AddCostAmountSum:   in.AddCostAmountSum,
		AddRefundAmountSum: in.AddRefundAmountSum,
		AddPaidOrderCnt:    in.AddPaidOrderCnt,
		AddRefundOrderCnt:  in.AddRefundOrderCnt,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserStatResp{}, nil
}
