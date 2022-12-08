package logic

import (
	"context"
	"fmt"
	"github.com/silenceper/wechat/v2/miniprogram/urllink"
	"github.com/silenceper/wechat/v2/miniprogram/urlscheme"
	"jakarta/common/xerr"

	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenWxMpUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenWxMpUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenWxMpUrlLogic {
	return &GenWxMpUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  生成小程序url link/schema
func (l *GenWxMpUrlLogic) GenWxMpUrl(in *pb.GenWxMpUrlReq) (*pb.GenWxMpUrlResp, error) {
	// todo: add your logic here and delete this line
	switch in.Type {
	case 1:
		return l.getUrlLink(in)
	case 2:
		return l.getUrlSchema(in)
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "平台参数错误")

	}
}

func (l *GenWxMpUrlLogic) getUrlLink(in *pb.GenWxMpUrlReq) (*pb.GenWxMpUrlResp, error) {
	if l.svcCtx.MiniProgram == nil {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "小程序未初始化")
	}
	parm := urllink.ULParams{
		Path:           in.Path,
		Query:          in.Query,
		EnvVersion:     "release",
		IsExpire:       true,
		ExpireType:     urllink.ExpireTypeInterval,
		ExpireTime:     0,
		ExpireInterval: int(in.ExpireIntervalDays),
	}
	resp := pb.GenWxMpUrlResp{}
	var err error
	resp.Url, err = l.svcCtx.MiniProgram.GetURLLink().Generate(&parm)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%+v", err))
	}
	return &resp, nil
}

func (l *GenWxMpUrlLogic) getUrlSchema(in *pb.GenWxMpUrlReq) (*pb.GenWxMpUrlResp, error) {
	if l.svcCtx.MiniProgram == nil {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "小程序未初始化")
	}
	parm := urlscheme.USParams{
		JumpWxa: &urlscheme.JumpWxa{
			Path:       in.Path,
			Query:      in.Query,
			EnvVersion: "release",
		},
		ExpireType:     urlscheme.ExpireTypeInterval,
		ExpireTime:     0,
		ExpireInterval: int(in.ExpireIntervalDays),
	}
	resp := pb.GenWxMpUrlResp{}
	var err error
	resp.Url, err = l.svcCtx.MiniProgram.GetSURLScheme().Generate(&parm)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%+v", err))
	}
	return &resp, nil
}
