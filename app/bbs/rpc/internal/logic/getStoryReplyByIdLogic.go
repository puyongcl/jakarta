package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/tool"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryReplyByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoryReplyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryReplyByIdLogic {
	return &GetStoryReplyByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XX的回复（根据id）
func (l *GetStoryReplyByIdLogic) GetStoryReplyById(in *pb.GetStoryReplyByIdReq) (*pb.GetStoryReplyByIdResp, error) {
	rs, err := l.svcCtx.StoryReplyModel.FindOne(l.ctx, in.StoryReplyId)
	if err != nil && err != bbsPgModel.ErrNotFound {
		return nil, err
	}

	resp := pb.GetStoryReplyByIdResp{Reply: &pb.StoryReply{}}
	if rs != nil {
		_ = copier.Copy(resp.Reply, rs)
		resp.Reply.CreateTime = tool.GetTimeDurationText(time.Now().Sub(rs.CreateTime)) + "前"
	}

	return &resp, nil
}
