package chat

import (
	"context"
	"encoding/json"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/kqueue"
	"jakarta/common/third_party/tim"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTextMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendTextMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTextMsgLogic {
	return &SendTextMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendTextMsgLogic) SendTextMsg(req *types.SendTextMsgReq) (resp *types.SendTextMsgResp, err error) {
	for idx := 0; idx < len(req.Text); idx++ {
		kqmsg := kqueue.SendImDefineMessage{
			FromUid: req.FromUid,
			ToUid:   req.ToUid,
			MsgType: req.MsgType,
			Text:    req.Text[idx],
			Sync:    tim.TimMsgSyncFromYes,
		}

		var buf []byte
		buf, err = json.Marshal(&kqmsg)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SendTextMsgLogic json err:%+v", err)
			return nil, err
		}

		err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("SendTextMsgLogic Push err:%+v", err)
			return nil, err
		}
	}

	return
}
