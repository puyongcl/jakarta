package admin

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/admin/api/internal/logic/adminlog"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterAdminReq) (resp *types.RegisterAdminResp, err error) {
	if req.AuthType != userkey.UserAuthTypePasswd || req.AuthKey == "" || req.Password == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "param is empty")
	}
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "Register", req.AdminUid, err, req, resp)
	}()
	var in pb.RegisterReq
	_ = copier.Copy(&in, req)
	in.UserType = userkey.UserTypeAdmin
	rsp, err := l.svcCtx.UsercenterRpc.Register(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.RegisterAdminResp{}
	_ = copier.Copy(resp, rsp.User)
	return
}
