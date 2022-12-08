package logic

import (
	"context"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserTypeLogic {
	return &UpdateUserTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserTypeLogic) UpdateUserType(in *pb.UpdateUserTypeReq) (*pb.UpdateUserTypeResp, error) {
	err := l.svcCtx.UserAuthModel.UpdateUserType(l.ctx, in.Uid, in.UserType)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserTypeResp{}, nil
}
