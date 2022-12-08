package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/bbs/rpc/internal/svc"
	"jakarta/app/bbs/rpc/pb"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/common/key/bbskey"
	"jakarta/common/key/db"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStoryReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddStoryReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStoryReplyLogic {
	return &AddStoryReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  回复
func (l *AddStoryReplyLogic) AddStoryReply(in *pb.AddStoryReplyReq) (*pb.AddStoryReplyResp, error) {
	// 查询是否重复
	var textMd5 string
	if in.ReplyText != "" {
		textMd5 = tool.Md5ByString(in.ReplyText)
	} else if in.ReplyVoice != "" {
		textMd5 = tool.Md5ByString(in.ReplyVoice)
	} else {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "内容为空")
	}

	cnt, err := l.svcCtx.StoryReplyModel.CountByTextMd5(l.ctx, textMd5)
	if err != nil {
		return nil, err
	}

	if cnt > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "重复内容")
	}

	// 查询XX
	st, err := l.svcCtx.StoryModel.FindOne(l.ctx, in.StoryId)
	if err != nil {
		return nil, err
	}
	if st.State == bbskey.StoryStateDeleted {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryNotFound, "已经删除，不能回复")
	}
	storyReplyId := fmt.Sprintf(db.DBUidId4, in.ListenerUid, in.StoryId, time.Now().Unix())

	data := &bbsPgModel.StoryReply{
		Id:          storyReplyId,
		StoryId:     in.StoryId,
		ListenerUid: in.ListenerUid,
		ReplyText:   in.ReplyText,
		ReplyVoice:  in.ReplyVoice,
		LikeCnt:     0,
		State:       bbskey.StoryStatusNotCheck,
		Uid:         st.Uid,
		TextMd5:     textMd5,
	}

	_, err = l.svcCtx.StoryReplyModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	// update story reply cnt
	err = l.svcCtx.StoryModel.AddStoryReplyCnt(l.ctx, in.StoryId)
	if err != nil {
		return nil, err
	}

	// 通知
	l.notify(in.Nickname, st.Tittle, st.Id, data.Id, in.ReplyText, st.StoryType, st.Uid)

	return &pb.AddStoryReplyResp{Id: data.Id}, nil
}

//  回复 限制1次
func (l *AddStoryReplyLogic) AddStoryReply2(in *pb.AddStoryReplyReq) (*pb.AddStoryReplyResp, error) {
	// 查询XX
	st, err := l.svcCtx.StoryModel.FindOne(l.ctx, in.StoryId)
	if err != nil {
		return nil, err
	}
	if st.State == bbskey.StoryStateDeleted {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryNotFound, "已经删除，不能回复")
	}
	storyReplyId := fmt.Sprintf(db.DBUidId3, in.ListenerUid, in.StoryId)
	// 查询回复
	data, err := l.svcCtx.StoryReplyModel.FindOne(l.ctx, storyReplyId)
	if err != nil && err != bbsPgModel.ErrNotFound {
		return nil, err
	}

	if data == nil && err == bbsPgModel.ErrNotFound { // 首次回复
		data = &bbsPgModel.StoryReply{
			Id:          storyReplyId,
			StoryId:     in.StoryId,
			ListenerUid: in.ListenerUid,
			ReplyText:   in.ReplyText,
			ReplyVoice:  in.ReplyVoice,
			LikeCnt:     0,
			State:       bbskey.StoryStatusNotCheck,
			Uid:         st.Uid,
		}

		_, err = l.svcCtx.StoryReplyModel.Insert(l.ctx, data)
		if err != nil {
			return nil, err
		}

		// update story reply cnt
		err = l.svcCtx.StoryModel.AddStoryReplyCnt(l.ctx, in.StoryId)
		if err != nil {
			return nil, err
		}

		// 通知
		l.notify(in.Nickname, st.Tittle, st.Id, data.Id, in.ReplyText, st.StoryType, st.Uid)

		return &pb.AddStoryReplyResp{Id: data.Id}, nil
	}
	if data.State == bbskey.StoryStateDeleted { // 已经删除 则更新
		err = l.svcCtx.StoryReplyModel.UpdateStoryNewReply(l.ctx, data.Id, in.ReplyText, in.ReplyVoice)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorAlreadyReply, "只能回复一次")
	}

	return &pb.AddStoryReplyResp{Id: data.Id}, nil
}

func (l *AddStoryReplyLogic) notify(listenerNickname, title, storyId, storyReplyId, content string, tp, toUid int64) {
	// 发送消息
	mpMsg := kqueue.SendMiniProgramSubscribeMessage{
		Thing1:  tool.CutText(title, 20, "..."),
		Thing3:  tool.CutText(content, 20, "..."),
		Date4:   time.Now().Format(db.DateTimeFormat),
		MsgType: notify.DefineNotifyMsgTypeMiniProgramMsg2,
		ToUid:   toUid,
		Page:    fmt.Sprintf(notify.DefineNotifyMsgTypeMiniProgramMsg2Path, storyReplyId),
	}
	imMsg := kqueue.SendImDefineMessage{
		FromUid:           notify.TimSystemNotifyUid,
		ToUid:             toUid,
		MsgType:           notify.DefineNotifyMsgTypeSystemMsg28,
		Title:             fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsgTitle28, listenerNickname, bbskey.GetStoryTypeText(tp)),
		Text:              tool.CutText(content, 40, "..."),
		Val1:              storyId,
		Val2:              storyReplyId,
		Val3:              "",
		Val4:              "",
		Val5:              "",
		Val6:              "",
		Sync:              tim.TimMsgSyncFromNo,
		RepeatMsgCheckId:  "",
		RepeatMsgCheckSec: 0,
	}
	kqMsg := kqueue.SubscribeNotifyMsgMessage{
		Uid:       toUid,
		TargetUid: 0,
		MsgType:   notify.DefineNotifyMsgTypeMiniProgramMsg2,
		SendCnt:   0,
		Action:    notify.SubscribeOneTimeNotifyMsgEventSend,
		IMMsg:     &imMsg,
		MpMsg:     &mpMsg,
	}
	buf, err := json.Marshal(kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddStory notify json marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueSendSubscribeNotifyMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddStory notify kafka Push err:%+v", err)
		return
	}
}
