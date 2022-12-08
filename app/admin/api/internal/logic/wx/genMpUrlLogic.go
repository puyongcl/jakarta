package wx

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/service"
	pbIm "jakarta/app/im/rpc/pb"
	"jakarta/common/xerr"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenMpUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenMpUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenMpUrlLogic {
	return &GenMpUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenMpUrlLogic) GenMpUrl(req *types.GenWxMpUrlReq) (resp *types.GenWxMpUrlResp, err error) {
	if l.svcCtx.Config.Mode != service.ProMode {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "非正式环境不支持")
	}
	var in pbIm.GenWxMpUrlReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.ImRpc.GenWxMpUrl(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GenWxMpUrlResp{}
	_ = copier.Copy(resp, rsp)

	return
}
