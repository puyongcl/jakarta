package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strconv"
	"strings"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedbackOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFeedbackOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedbackOrderLogic {
	return &FeedbackOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  XXX反馈
func (l *FeedbackOrderLogic) FeedbackOrder(in *pb.FeedbackOrderReq) (*pb.FeedbackOrderResp, error) {
	orderData, err := l.svcCtx.ChatOrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	if orderData.Feedback != "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.OrderErrorAlreadyFeedback, "已经反馈，无法修改")
	}
	newData := new(orderPgModel.ChatOrder)
	newData.Feedback = in.Feedback
	newData.OrderId = in.OrderId
	err = l.svcCtx.ChatOrderModel.UpdateOrderOpinionAndState(l.ctx, in.OrderId, newData, orderkey.NeedFeedbackOderState)
	if err != nil {
		return nil, err
	}

	// 通知消息
	kqMsgs := make([]*kqueue.SendImDefineMessage, 0)
	kqMsgs = append(kqMsgs, &kqueue.SendImDefineMessage{
		FromUid: notify.TimOrderNotifyUid,
		ToUid:   in.Uid,
		MsgType: notify.DefineNotifyMsgTypeOrderMsg6,
		Title:   notify.DefineNotifyMsgTemplateOrderMsgTitle6,
		Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateOrderMsg6, strings.Replace(orderData.ListenerNickName, "#", "", -1)),
		Val1:    strconv.FormatInt(in.ListenerUid, 10),
		Val2:    orderData.OrderId,
		Sync:    tim.TimMsgSyncFromNo,
	})
	// 同步反馈到聊天
	if in.SendMsg == db.Enable && in.Feedback != "" {
		// 对用户
		msg1 := kqueue.SendImDefineMessage{
			FromUid: orderData.ListenerUid,
			ToUid:   in.Uid,
			MsgType: notify.DefineNotifyMsgTypeChatMsg25,
			Title:   "",
			Text:    in.Feedback,
			Val1:    orderData.OrderId,
			Val2:    "",
			Sync:    tim.TimMsgSyncFromYes,
		}
		kqMsgs = append(kqMsgs, &msg1)
	}
	// 对XXX
	msg2 := kqueue.SendImDefineMessage{
		FromUid: orderData.Uid,
		ToUid:   in.ListenerUid,
		MsgType: notify.DefineNotifyMsgTypeChatMsg27,
		Title:   "",
		Text:    notify.DefineNotifyMsgTemplateChatMsg27,
		Val1:    "",
		Val2:    "",
		Sync:    tim.TimMsgSyncFromNo,
	}
	kqMsgs = append(kqMsgs, &msg2)

	l.notify(kqMsgs)
	l.notifyMp(orderData.Uid, orderData.ListenerNickName, in.Feedback)
	return &pb.FeedbackOrderResp{}, nil
}

func (l *FeedbackOrderLogic) notify(kqMsgs []*kqueue.SendImDefineMessage) {
	for k, _ := range kqMsgs {
		buf, err := json.Marshal(kqMsgs[k])
		if err != nil {
			logx.WithContext(l.ctx).Errorf("FeedbackOrderLogic json marshal err:%+v", err)
			return
		}
		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("FeedbackOrderLogic Push err:%+v", err)
			return
		}
	}
	return
}

// 同步发送到小程序服务通知
func (l *FeedbackOrderLogic) notifyMp(uid int64, listenerNickname string, content string) {
	// 判断用户是否有订阅
	b, err := l.svcCtx.IMRedis.IsUserHaveOneTimeSubscribeMsg(l.ctx, notify.DefineNotifyMsgTypeMiniProgramMsg3, uid)
	if err != nil {
		return
	}
	if !b {
		return
	}

	// 发送消息
	mpMsg := kqueue.SendMiniProgramSubscribeMessage{
		Thing1:  tool.CutText(listenerNickname, 20, "..."),
		Thing3:  tool.CutText(content, 20, "..."),
		Time2:   time.Now().Format(db.DateTimeFormat),
		MsgType: notify.DefineNotifyMsgTypeMiniProgramMsg3,
		ToUid:   uid,
		Page:    notify.DefineNotifyMsgTypeMiniProgramMsg3Path,
	}

	kqMsg := kqueue.SubscribeNotifyMsgMessage{
		Uid:       uid,
		TargetUid: 0,
		MsgType:   notify.DefineNotifyMsgTypeMiniProgramMsg3,
		SendCnt:   0,
		Action:    notify.SubscribeOneTimeNotifyMsgEventSend,
		MpMsg:     &mpMsg,
	}
	buf, err := json.Marshal(kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddStory notify json marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqSendSubscribeNotifyMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddStory notify kafka Push err:%+v", err)
		return
	}
	return
}
