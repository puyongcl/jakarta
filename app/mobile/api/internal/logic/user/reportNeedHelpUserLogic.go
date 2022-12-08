package user

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/usercenter/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportNeedHelpUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportNeedHelpUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportNeedHelpUserLogic {
	return &ReportNeedHelpUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportNeedHelpUserLogic) ReportNeedHelpUser(req *types.ReportNeedHelpUserReq) (resp *types.ReportNeedHelpUserResp, err error) {
	var in pb.ReportNeedHelpUserReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.UsercenterRpc.ReportNeedHelpUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.ReportNeedHelpUserResp{}
	return
}
