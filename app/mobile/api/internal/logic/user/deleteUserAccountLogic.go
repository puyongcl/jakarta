package user

import (
	"context"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/userkey"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserAccountLogic {
	return &DeleteUserAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserAccountLogic) DeleteUserAccount(req *types.DeleteUserAccountReq) (resp *types.DeleteUserAccountResp, err error) {
	rsp, err := l.svcCtx.UsercenterRpc.DeleteUserAccount(l.ctx, &pb.DeleteUserAccountReq{Uid: req.Uid})
	if err != nil {
		return nil, err
	}
	if rsp.UserType == userkey.UserTypeListener {
		_, err = l.svcCtx.ListenerRpc.ChangeWorkState(l.ctx, &pbListener.ChangeWorkStateReq{
			ListenerUid: req.Uid,
			WorkState:   listenerkey.ListenerWorkStateAccountDeleted,
		})
		if err != nil {
			return nil, err
		}
	}
	resp = &types.DeleteUserAccountResp{}
	return
}
