package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryReplyListByOwnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoryReplyListByOwnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryReplyListByOwnLogic {
	return &GetStoryReplyListByOwnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取个人所有回复
func (l *GetStoryReplyListByOwnLogic) GetStoryReplyListByOwn(in *pb.GetStoryReplyListByOwnReq) (*pb.GetStoryReplyListByOwnResp, error) {
	rs, err := l.svcCtx.StoryReplyModel.Find(l.ctx, in.ListenerUid, "", in.PageNo, in.PageSize, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.GetStoryReplyListByOwnResp{List: make([]*pb.StoryReply, 0)}
	if len(rs) > 0 {
		now := time.Now()
		for idx := 0; idx < len(rs); idx++ {
			var val pb.StoryReply
			val = pb.StoryReply{
				CreateTime:  tool.GetTimeDurationText(now.Sub(rs[idx].CreateTime)) + "前",
				Id:          rs[idx].Id,
				StoryId:     rs[idx].StoryId,
				ListenerUid: rs[idx].ListenerUid,
				ReplyText:   rs[idx].ReplyText,
				ReplyVoice:  rs[idx].ReplyVoice,
				LikeCnt:     rs[idx].LikeCnt,
				State:       rs[idx].State,
				Uid:         rs[idx].Uid,
			}

			// 是否点赞
			var srll *bbsPgModel.StoryReplyLikeLog
			srll, err = l.svcCtx.StoryReplyLikeLogModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId3, in.ListenerUid, val.StoryId))
			if err != nil && err != bbsPgModel.ErrNotFound {
				return nil, err
			}
			err = nil
			if srll != nil {
				val.IsLike = db.Enable
			} else {
				val.IsLike = db.Disable
			}

			resp.List = append(resp.List, &val)
		}
	}

	return &resp, nil
}
