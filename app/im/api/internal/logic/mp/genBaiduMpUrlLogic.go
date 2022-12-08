package mp

import (
	"context"
	pbIm "jakarta/app/im/rpc/pb"

	"jakarta/app/im/api/internal/svc"
	"jakarta/app/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenBaiduMpUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenBaiduMpUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenBaiduMpUrlLogic {
	return &GenBaiduMpUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenBaiduMpUrlLogic) GenBaiduMpUrl(req *types.GenBaiduWxMpUrlReq) (resp *types.GenBaiduWxMpUrlResp) {
	var in pbIm.GenWxMpUrlReq
	in.Query = "source=baidu"
	in.Path = req.Path
	in.ExpireIntervalDays = 30
	in.Type = 2
	resp = &types.GenBaiduWxMpUrlResp{
		Code:   0,
		Msg:    "",
		Result: types.BaiduWxMpUrl{},
	}
	rsp, err := l.svcCtx.ImRpc.GenWxMpUrl(l.ctx, &in)
	if err != nil {
		resp.Code = 1
		return
	}
	resp.Result.Scheme = rsp.Url
	return
}
