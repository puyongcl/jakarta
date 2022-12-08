package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReportUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReportUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReportUserListLogic {
	return &GetReportUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取上报用户列表
func (l *GetReportUserListLogic) GetReportUserList(in *pb.GetReportUserListReq) (*pb.GetReportUserListResp, error) {
	cnt, err := l.svcCtx.UserReportModel.FindCount(l.ctx, in.Uid, in.TargetUid, in.Tag, in.State)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetReportUserListResp{List: make([]*pb.ReportUserData, 0), Sum: cnt}
	if cnt <= 0 {
		return resp, nil
	}
	rs, err := l.svcCtx.UserReportModel.Find(l.ctx, in.Uid, in.TargetUid, in.Tag, in.State, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}

	for idx := 0; idx < len(rs); idx++ {
		var val pb.ReportUserData
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].CreateTime.Format(db.DateTimeFormat)
		val.UpdateTime = rs[idx].UpdateTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
