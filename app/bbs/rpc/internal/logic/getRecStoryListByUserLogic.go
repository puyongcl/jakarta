package logic

import (
	"context"
	"jakarta/common/tool"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecStoryListByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecStoryListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecStoryListByUserLogic {
	return &GetRecStoryListByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  推荐
func (l *GetRecStoryListByUserLogic) GetRecStoryListByUser(in *pb.GetRecStoryListByUserReq) (*pb.GetRecStoryListByUserResp, error) {
	// 获取黑名单
	blk, err := l.svcCtx.UserRedis.GetBlacklist(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	rs, err := l.svcCtx.StoryModel.FindRec(l.ctx, in.StoryType, in.Spec, in.PageNo, in.PageSize, blk)
	if err != nil {
		return nil, err
	}
	resp := pb.GetRecStoryListByUserResp{List: make([]*pb.Story, 0)}

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
