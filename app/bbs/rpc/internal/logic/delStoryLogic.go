package logic

import (
	"context"
	"jakarta/common/key/bbskey"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelStoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelStoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelStoryLogic {
	return &DelStoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  删除发布
func (l *DelStoryLogic) DelStory(in *pb.DelStoryReq) (*pb.DelStoryResp, error) {
	err := l.svcCtx.StoryModel.UpdateStoryState(l.ctx, in.StoryId, bbskey.StoryStateDeleted)
	if err != nil {
		return nil, err
	}
	return &pb.DelStoryResp{}, nil
}
