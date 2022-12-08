package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	"jakarta/app/im/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/xerr"
	"time"
)

// 发送自定义IM消息
type SendDefineImMsgMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendDefineImMsgMq(ctx context.Context, svcCtx *svc.ServiceContext) *SendDefineImMsgMq {
	return &SendDefineImMsgMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendDefineImMsgMq) Consume(_, val string) error {
	var message kqueue.SendImDefineMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendDefineImMsgMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendDefineImMsgMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *SendDefineImMsgMq) execService(message *kqueue.SendImDefineMessage) error {
	// 限定频率和去重
	b, err := l.svcCtx.IMRedis.SetNotifyLimit(l.ctx, message.RepeatMsgCheckId, message.MsgType, message.RepeatMsgCheckSec)
	if err != nil {
		return err
	}
	if !b {
		return nil
	}

	// 某些消息需要同时发送到服务号
	switch message.MsgType {
	case notify.DefineNotifyMsgTypeViewMsg17: // 浏览通知
		//in := pbUser.GetUserOnlineStateReq{Uid: message.ToUid}
		//var rsp *pbUser.GetUserOnlineStateResp
		//rsp, err = l.svcCtx.UserRpc.GetUserOnlineState(l.ctx, &in)
		//if err != nil {
		//	return err
		//}
		//if rsp.OnlineState == userkey.Login {
		//	break
		//}
		err = l.syncSendFwh(message)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SendDefineImMsgMq err:%+v", err)
		}

	case notify.DefineNotifyMsgTypeOrderMsg1, notify.DefineNotifyMsgTypeOrderMsg2,
		notify.DefineNotifyMsgTypeSystemMsg19, notify.DefineNotifyMsgTypeSystemMsg20, notify.DefineNotifyMsgTypeSystemMsg21,
		notify.DefineNotifyMsgTypeSystemMsg22:
		err = l.syncSendFwh(message)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SendDefineImMsgMq err:%+v", err)
		}

	default:

	}

	// 发送im消息
	switch message.MsgType {
	case notify.TextMsgTypeIntroMsg1, notify.TextMsgTypeHelloMsg2, notify.TextMsgTypeAdviserMsg3: // 普通文本消息
		var in pb.SendTextMsgReq
		_ = copier.Copy(&in, message)
		_, err = l.svcCtx.ImRpc.SendTextMsg(l.ctx, &in)
		if err != nil {
			return err
		}
	default: // 自定义消息
		var in pb.SendDefineMsgReq
		_ = copier.Copy(&in, message)
		_, err = l.svcCtx.ImRpc.SendDefineMsg(l.ctx, &in)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *SendDefineImMsgMq) syncSendFwh(kqMsg *kqueue.SendImDefineMessage) error {
	fwhMsg := new(kqueue.SendFwhNotifyMessage)
	switch kqMsg.MsgType {
	case notify.DefineNotifyMsgTypeViewMsg17:
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg1FirstData,
			Keyword1: kqMsg.Val3,
			Keyword2: time.Now().Format(db.DateTimeFormat),
			Keyword3: "",
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg1RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg1,
			ToUid:    kqMsg.ToUid,
		}
	case notify.DefineNotifyMsgTypeOrderMsg1, notify.DefineNotifyMsgTypeOrderMsg2:
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg2FirstData,
			Keyword1: kqMsg.Val3,
			Keyword2: kqMsg.Val4,
			Keyword3: kqMsg.Val5,
			Keyword4: kqMsg.Val6,
			Remark:   notify.DefineNotifyMsgTypeFwhMsg2RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg2,
			ToUid:    kqMsg.ToUid,
			Color:    notify.DefineNotifyMsgTypeFwhMsg2Color,
		}

	case notify.DefineNotifyMsgTypeSystemMsg19: // XXX审核不通过
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg3FirstData2,
			Keyword1: kqMsg.Val3,
			Keyword2: kqMsg.Val4,
			Keyword3: kqMsg.Val5,
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg3RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg3,
			ToUid:    kqMsg.ToUid,
		}
	case notify.DefineNotifyMsgTypeSystemMsg20: // XXX审核通过
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg3FirstData1,
			Keyword1: kqMsg.Val3,
			Keyword2: kqMsg.Val4,
			Keyword3: kqMsg.Val5,
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg3RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg3,
			ToUid:    kqMsg.ToUid,
		}
	case notify.DefineNotifyMsgTypeSystemMsg21: // 修改资料审核通过
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg3FirstData3,
			Keyword1: kqMsg.Val3,
			Keyword2: kqMsg.Val4,
			Keyword3: kqMsg.Val5,
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg3RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg3,
			ToUid:    kqMsg.ToUid,
		}
	case notify.DefineNotifyMsgTypeSystemMsg22: // 修改资料审核不通过
		*fwhMsg = kqueue.SendFwhNotifyMessage{
			First:    notify.DefineNotifyMsgTypeFwhMsg3FirstData4,
			Keyword1: kqMsg.Val3,
			Keyword2: kqMsg.Val4,
			Keyword3: kqMsg.Val5,
			Keyword4: "",
			Remark:   notify.DefineNotifyMsgTypeFwhMsg3RemarkData,
			MsgType:  notify.DefineNotifyMsgTypeFwhMsg3,
			ToUid:    kqMsg.ToUid,
		}
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数错误 %d", kqMsg.MsgType))
	}

	//
	buf, err := json.Marshal(fwhMsg)
	if err != nil {
		return err
	}
	err = l.svcCtx.KqueueSendWxFwhMsgClient.Push(string(buf))
	if err != nil {
		return err
	}
	return nil
}
