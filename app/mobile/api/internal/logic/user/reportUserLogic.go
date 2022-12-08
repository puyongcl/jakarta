package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportUserLogic {
	return &ReportUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportUserLogic) ReportUser(req *types.ReportUserReq) (resp *types.ReportUserResp, err error) {
	var in pb.ReportUserReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.UsercenterRpc.ReportUser(l.ctx, &in)
	if err != nil {

	}
	resp = &types.ReportUserResp{}
	return
}
