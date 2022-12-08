package logic

import (
	"context"
	"jakarta/common/xerr"

	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserUnionIdByFwhOpenIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserUnionIdByFwhOpenIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserUnionIdByFwhOpenIdLogic {
	return &GetUserUnionIdByFwhOpenIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  根据用户的服务号openid获取用户的unionId
func (l *GetUserUnionIdByFwhOpenIdLogic) GetUserUnionIdByFwhOpenId(in *pb.GetUserUnionIdByFwhOpenIdReq) (*pb.GetUserUnionIdByFwhOpenIdResp, error) {
	if in.OpenId == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	uio, err := l.svcCtx.Wxfwh.GetUser().GetUserInfo(in.OpenId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserUnionIdByFwhOpenIdResp{UnionId: uio.UnionID}, nil
}
