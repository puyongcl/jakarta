package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/bbs/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryByIdLogic {
	return &GetStoryByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStoryByIdLogic) GetStoryById(req *types.GetStoryByIdReq) (resp *types.GetStoryByIdResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	var in pb.GetStoryByIdReq
	_ = copier.Copy(&in, req)
	in.Uid = uid
	rsp, err := l.svcCtx.BbsRpc.GetStoryById(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	var val types.Story
	_ = copier.Copy(&val, rsp.Story)

	var in2 pbUser.GetUserShortProfileReq
	in2.Uid = rsp.Story.Uid
	rsp2, err := l.svcCtx.UsercenterRpc.GetUserShortProfile(l.ctx, &in2)
	if err != nil {
		return nil, err
	}
	val.Avatar = rsp2.User.Avatar
	val.Nickname = rsp2.User.Nickname

	resp = &types.GetStoryByIdResp{Story: &val}
	return
}
