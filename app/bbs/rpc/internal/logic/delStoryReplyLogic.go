package logic

import (
	"context"
	"jakarta/common/key/bbskey"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStoryReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelStoryReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStoryReplyLogic {
	return &DelStoryReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  删除回复
func (l *DelStoryReplyLogic) DelStoryReply(in *pb.DelStoryReplyReq) (*pb.DelStoryReplyResp, error) {
	err := l.svcCtx.StoryReplyModel.UpdateStoryReplyState(l.ctx, in.StoryReplyId, bbskey.StoryStateDeleted)
	if err != nil {
		return nil, err
	}
	return &pb.DelStoryReplyResp{}, nil
}
