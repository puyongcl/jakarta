package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/bbs/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
)

type GetStoryReplyByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStoryReplyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryReplyByIdLogic {
	return &GetStoryReplyByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStoryReplyByIdLogic) GetStoryReplyById(req *types.GetStoryReplyReq) (resp *types.GetStoryReplyResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	var in pb.GetStoryReplyByIdReq
	_ = copier.Copy(&in, req)
	in.Uid = uid
	rsp, err := l.svcCtx.BbsRpc.GetStoryReplyById(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetStoryReplyResp{Reply: &types.StoryReply{}}
	_ = copier.Copy(resp.Reply, rsp.Reply)

	var lin pbListener.GetListenerBasicInfoReq
	lin.ListenerUid = rsp.Reply.ListenerUid
	rsp2, err := l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &lin)
	if err != nil {
		return nil, err
	}
	resp.Reply.ListenerNickname = rsp2.NickName
	resp.Reply.ListenerAvatar = rsp2.Avatar

	return
}
