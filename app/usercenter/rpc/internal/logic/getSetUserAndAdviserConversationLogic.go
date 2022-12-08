package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSetUserAndAdviserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSetUserAndAdviserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSetUserAndAdviserConversationLogic {
	return &GetSetUserAndAdviserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取和更新用户与顾问的对话
func (l *GetSetUserAndAdviserConversationLogic) GetSetUserAndAdviserConversation(in *pb.GetSetUserAndAdviserConversationReq) (*pb.GetSetUserAndAdviserConversationResp, error) {
	var uac *userPgModel.UserAdviserConversation
	var err error
	uac, err = l.svcCtx.UserAdviserConversationModel.FindOne(l.ctx, in.Uid)
	if err != nil && err != userPgModel.ErrNotFound {
		return nil, err
	}
	resp := pb.GetSetUserAndAdviserConversationResp{}
	if uac != nil {
		var up int64
		if in.Step > uac.Step {
			uac.Step = in.Step
			up++
		}
		if len(in.SelectSpec) > 0 && !tool.IsEqualArrayInt64(in.SelectSpec, uac.SelectSpec) {
			uac.SelectSpec = in.SelectSpec
			up++
		}
		if len(in.Conversation) > len(uac.Conversation) {
			uac.Conversation = in.Conversation
			up++
		}

		if up > 0 {
			err = l.svcCtx.UserAdviserConversationModel.Update(l.ctx, uac)
			if err != nil {
				return nil, err
			}
		}

		_ = copier.Copy(&resp, uac)
	} else {
		if in.Step <= 0 {
			return &resp, nil
		}
		uac = &userPgModel.UserAdviserConversation{
			Uid:          in.Uid,
			Step:         in.Step,
			SelectSpec:   in.SelectSpec,
			Conversation: in.Conversation,
		}
		_, err = l.svcCtx.UserAdviserConversationModel.Insert(l.ctx, uac)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&resp, uac)
	}
	return &resp, nil
}
