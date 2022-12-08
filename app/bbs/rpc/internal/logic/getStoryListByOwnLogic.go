package logic

import (
	"context"
	"jakarta/common/tool"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryListByOwnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoryListByOwnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryListByOwnLogic {
	return &GetStoryListByOwnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取个人所有XX
func (l *GetStoryListByOwnLogic) GetStoryListByOwn(in *pb.GetStoryListByOwnReq) (*pb.GetStoryListByOwnResp, error) {
	rs, err := l.svcCtx.StoryModel.Find(l.ctx, in.Uid, in.StoryType, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := pb.GetStoryListByOwnResp{List: make([]*pb.Story, 0)}
	if len(rs) > 0 {
		now := time.Now()
		for idx := 0; idx < len(rs); idx++ {
			var val pb.Story
			val = pb.Story{
				CreateTime: tool.GetTimeDurationText(now.Sub(rs[idx].CreateTime)) + "前",
				Id:         rs[idx].Id,
				Uid:        rs[idx].Uid,
				StoryType:  rs[idx].StoryType,
				Spec:       rs[idx].Spec,
				Tittle:     rs[idx].Tittle,
				Content:    rs[idx].Content,
				State:      rs[idx].State,
				ViewCnt:    rs[idx].ViewCnt,
				ReplyCnt:   rs[idx].ReplyCnt,
			}

			resp.List = append(resp.List, &val)
		}
	}

	return &resp, nil
}
