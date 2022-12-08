package user

import (
	"context"
	"github.com/jinzhu/copier"
	pbUser "jakarta/app/usercenter/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSetUserAndAdviserConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSetUserAndAdviserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSetUserAndAdviserConversationLogic {
	return &GetSetUserAndAdviserConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSetUserAndAdviserConversationLogic) GetSetUserAndAdviserConversation(req *types.GetSetUserAndAdviserConversationReq) (resp *types.GetSetUserAndAdviserConversationResp, err error) {
	var in pbUser.GetSetUserAndAdviserConversationReq
	in.Uid = req.Uid
	in.Step = req.Step
	in.Conversation = req.Conversation
	in.SelectSpec = req.SelectSpec
	rsp, err := l.svcCtx.UsercenterRpc.GetSetUserAndAdviserConversation(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetSetUserAndAdviserConversationResp{}
	_ = copier.Copy(resp, rsp)
	return
}
