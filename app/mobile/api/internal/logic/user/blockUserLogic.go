package user

import (
	"context"
	"github.com/jinzhu/copier"
	pbUser "jakarta/app/usercenter/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlockUserLogic) BlockUser(req *types.BlockUserReq) (resp *types.BlockUserResp, err error) {
	var in pbUser.BlockUserReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.UsercenterRpc.BlockUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.BlockUserResp{}
	return
}
