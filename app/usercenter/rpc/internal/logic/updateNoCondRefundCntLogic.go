package logic

import (
	"context"
	"jakarta/common/key/db"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoCondRefundCntLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNoCondRefundCntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoCondRefundCntLogic {
	return &UpdateNoCondRefundCntLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNoCondRefundCntLogic) UpdateNoCondRefundCnt(in *pb.UpdateNoCondRefundCntReq) (*pb.UpdateNoCondRefundCntResp, error) {
	b, err := l.svcCtx.UserStatModel.CanNoCondRefund(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateNoCondRefundCntResp{State: db.Disable}
	if b {
		resp.State = db.Enable
	}
	return resp, nil
}
