package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/key/bbskey"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryByIdLogic {
	return &GetStoryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取发布的XX（根据id）
func (l *GetStoryByIdLogic) GetStoryById(in *pb.GetStoryByIdReq) (*pb.GetStoryByIdResp, error) {
	data, err := l.svcCtx.StoryModel.FindOne(l.ctx, in.StoryId)
	if err != nil && err != bbsPgModel.ErrNotFound {
		return nil, err
	}

	if data == nil {
		return &pb.GetStoryByIdResp{}, nil
	}
	if data.State == bbskey.StoryStateDeleted {
		return &pb.GetStoryByIdResp{Story: nil}, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryNotFound, "已经删除")
	}

	var val pb.Story
	now := time.Now()
	val = pb.Story{
		CreateTime: tool.GetTimeDurationText(now.Sub(data.CreateTime)) + "前",
		Id:         data.Id,
		Uid:        data.Uid,
		StoryType:  data.StoryType,
		Spec:       data.Spec,
		Tittle:     data.Tittle,
		Content:    data.Content,
		State:      data.State,
		ViewCnt:    data.ViewCnt,
		ReplyCnt:   data.ReplyCnt,
	}

	// 更新浏览次数
	storyViewLogId := fmt.Sprintf(db.DBUidId3, in.Uid, in.StoryId)
	svl, err := l.svcCtx.StoryViewLogModel.FindOne(l.ctx, storyViewLogId)
	if err != nil && err != bbsPgModel.ErrNotFound {
		return nil, err
	}

	if svl == nil && err == bbsPgModel.ErrNotFound { // 首次浏览记录
		svl = &bbsPgModel.StoryViewLog{
			Id:      storyViewLogId,
			ViewCnt: 1,
			Uid:     in.Uid,
			StoryId: in.StoryId,
		}
		_, err = l.svcCtx.StoryViewLogModel.Insert(l.ctx, svl)
		if err != nil {
			return nil, err
		}
		// story view cnt +1
		err = l.svcCtx.StoryModel.AddStoryViewCnt(l.ctx, in.StoryId)
		if err != nil {
			return nil, err
		}
	} else {
		err = l.svcCtx.StoryViewLogModel.AddStoryReplyViewLogLikeCnt(l.ctx, storyViewLogId)
		if err != nil {
			return nil, err
		}
		// story view cnt +1
		err = l.svcCtx.StoryModel.AddStoryViewCnt(l.ctx, in.StoryId)
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetStoryByIdResp{Story: &val}, nil
}
