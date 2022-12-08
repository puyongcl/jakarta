package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/key/bbskey"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddStoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStoryLogic {
	return &AddStoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发布
func (l *AddStoryLogic) AddStory(in *pb.AddStoryReq) (*pb.AddStoryResp, error) {
	// TODO 限制发布频率

	// 查询是否重复
	textMd5 := tool.Md5ByString(in.Content)
	cnt, err := l.svcCtx.StoryModel.CountByTextMd5(l.ctx, textMd5)
	if err != nil {
		return nil, err
	}

	if cnt > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "重复内容")
	}

	data := bbsPgModel.Story{
		Id:        uniqueid.GenDataId(),
		Uid:       in.Uid,
		StoryType: in.StoryType,
		Spec:      in.Spec,
		Tittle:    in.Tittle,
		Content:   in.Content,
		State:     bbskey.StoryStatusNotCheck,
		ViewCnt:   0,
		ReplyCnt:  0,
		TextMd5:   textMd5,
	}

	_, err = l.svcCtx.StoryModel.Insert(l.ctx, &data)
	if err != nil {
		return nil, err
	}

	// 通知XXX回复
	var listenerUids []int64
	listenerUids, err = l.svcCtx.ListenerRedis.GetRecommendListenerAndIncrScore(l.ctx, rediskey.RedisKeyListenerRecommendReplyUserStory, 3)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddStory ListenerRedis.GetRecommendListenerAndIncrScore err:%+v", err)
		err = nil
	}
	if len(listenerUids) > 0 {
		l.notify(in, listenerUids, data.Id)
	}
	return &pb.AddStoryResp{Id: data.Id}, nil
}

func (l *AddStoryLogic) notify(in *pb.AddStoryReq, listenerUid []int64, storyId string) {
	for idx := 0; idx < len(listenerUid); idx++ {
		fwhMsg := kqueue.SendFwhNotifyMessage{
			First:    fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg4FirstData, bbskey.GetStoryTypeText(in.StoryType)),
			Keyword1: in.Nickname,
			Keyword2: bbskey.GetStoryTypeText(in.StoryType),
			Keyword3: time.Now().Format(db.DateTimeFormat),
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg4RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg4,
			ToUid:    listenerUid[idx],
			Path:     fmt.Sprintf(notify.DefineNotifyMsgTypeFwhMsg4Path, storyId),
		}

		buf, err := json.Marshal(&fwhMsg)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("AddStory notify json marshal err:%+v", err)
			return
		}
		err = l.svcCtx.KqueueSendWxFwhMsgClient.Push(string(buf))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("AddStory notify kafka Push err:%+v", err)
			return
		}
	}

}
