package logic

import (
	"context"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWxOpenIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserWxOpenIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWxOpenIdLogic {
	return &GetUserWxOpenIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  根据uid获取用户的openid
func (l *GetUserWxOpenIdLogic) GetUserWxOpenId(in *pb.GetUserWxOpenIdReq) (*pb.GetUserWxOpenIdResp, error) {
	rs, err := l.svcCtx.UserWechatInfoModel.FindOneByUid(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserWxOpenIdResp{FwhOpenId: rs.FwhOpenid, MpOpenId: rs.MpOpenid, UnionId: rs.UnionId}, nil
}
