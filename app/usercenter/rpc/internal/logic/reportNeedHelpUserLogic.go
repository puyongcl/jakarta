package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/common/uniqueid"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportNeedHelpUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportNeedHelpUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportNeedHelpUserLogic {
	return &ReportNeedHelpUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  上报需要帮助的用户
func (l *ReportNeedHelpUserLogic) ReportNeedHelpUser(in *pb.ReportNeedHelpUserReq) (*pb.ReportNeedHelpUserResp, error) {
	data := new(userPgModel.UserNeedHelp)
	_ = copier.Copy(data, in)
	data.Id = uniqueid.GenDataId()
	_, err := l.svcCtx.UserNeedHelpModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}
	return &pb.ReportNeedHelpUserResp{}, nil
}
