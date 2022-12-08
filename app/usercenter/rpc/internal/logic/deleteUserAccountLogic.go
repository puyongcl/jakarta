package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/service"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserAccountLogic {
	return &DeleteUserAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  删除用户账户
func (l *DeleteUserAccountLogic) DeleteUserAccount(in *pb.DeleteUserAccountReq) (*pb.DeleteUserAccountResp, error) {
	data, err := l.svcCtx.UserAuthModel.DeleteUserAccount(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	resp := pb.DeleteUserAccountResp{UserType: data.UserType}
	if l.svcCtx.Config.Mode == service.DevMode {
		err = l.svcCtx.TimClient.AccountDel(in.Uid)
		if err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
