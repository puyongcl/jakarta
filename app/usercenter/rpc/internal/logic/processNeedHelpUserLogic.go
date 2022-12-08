package logic

import (
	"context"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessNeedHelpUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessNeedHelpUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessNeedHelpUserLogic {
	return &ProcessNeedHelpUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理记录需要帮助的用户处理结果
func (l *ProcessNeedHelpUserLogic) ProcessNeedHelpUser(in *pb.AdminMarkNeedHelpUserReq) (*pb.AdminMarkNeedHelpUserResp, error) {
	// TODO 处理
	// 更新状态
	err := l.svcCtx.UserNeedHelpModel.UpdateState(l.ctx, in.Id, in.Action, in.Remark)
	if err != nil {
		return nil, err
	}
	return &pb.AdminMarkNeedHelpUserResp{}, nil
}
