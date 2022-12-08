package user

import (
	"context"
	"jakarta/app/admin/api/internal/logic/adminlog"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/userkey"
	"time"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

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

var limitTime = time.Date(2022, 9, 30, 23, 59, 59, 0, time.Local)

func (l *DeleteUserAccountLogic) DeleteUserAccount(req *types.AdminDeleteUserAccountReq) (resp *types.AdminDeleteUserAccountResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "DeleteUserAccountLogic", req.AdminUid, err, req, resp)
	}()
	//if l.svcCtx.Config.Mode != service.DevMode && time.Now().Unix() > limitTime.Unix() {
	//	return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "不支持")
	//}
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
	resp = &types.AdminDeleteUserAccountResp{}

	return
}
