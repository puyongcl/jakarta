package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/uniqueid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportUserLogic {
	return &ReportUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  上报用户
func (l *ReportUserLogic) ReportUser(in *pb.ReportUserReq) (*pb.ReportUserResp, error) {
	data := new(userPgModel.UserReport)
	_ = copier.Copy(data, in)
	data.State = userkey.ReportStateCreated
	data.Id = uniqueid.GenDataId()
	_, err := l.svcCtx.UserReportModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}
	return &pb.ReportUserResp{}, nil
}
