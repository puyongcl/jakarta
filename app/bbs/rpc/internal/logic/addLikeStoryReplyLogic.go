package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/key/db"
)

type AddLikeStoryReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLikeStoryReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLikeStoryReplyLogic {
	return &AddLikeStoryReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  点赞
func (l *AddLikeStoryReplyLogic) AddLikeStoryReply(in *pb.AddLikeStoryReplyReq) (*pb.AddLikeStoryReplyResp, error) {
	// 查询是否存在
	likeLogId := fmt.Sprintf(db.DBUidId3, in.Uid, in.StoryReplyId)
	data, err := l.svcCtx.StoryReplyLikeLogModel.FindOne(l.ctx, likeLogId)
	if err != nil && err != bbsPgModel.ErrNotFound {
		return nil, err
	}
	if data == nil && err == bbsPgModel.ErrNotFound { // 第一次点赞
		data = &bbsPgModel.StoryReplyLikeLog{
			Id:           likeLogId,
			LikeCnt:      1,
			Uid:          in.Uid,
			StoryReplyId: in.StoryReplyId,
		}
		_, err = l.svcCtx.StoryReplyLikeLogModel.Insert(l.ctx, data)
		if err != nil {
			return nil, err
		}

		// story reply like cnt +1
		err = l.svcCtx.StoryReplyModel.AddStoryReplyLikeCnt(l.ctx, in.StoryReplyId)
		if err != nil {
			return nil, err
		}
		data.LikeCnt++
	} else {
		err = l.svcCtx.StoryReplyLikeLogModel.AddStoryReplyLikeLogLikeCnt(l.ctx, likeLogId)
		if err != nil {
			return nil, err
		}
	}
	return &pb.AddLikeStoryReplyResp{LikeCnt: data.LikeCnt}, nil
}
