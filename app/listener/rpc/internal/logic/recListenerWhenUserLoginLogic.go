package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/cservice"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecListenerWhenUserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecListenerWhenUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecListenerWhenUserLoginLogic {
	return &RecListenerWhenUserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取新用户推荐XXX
func (l *RecListenerWhenUserLoginLogic) RecListenerWhenUserLogin(in *pb.RecListenerWhenUserLoginReq) (*pb.RecListenerWhenUserLoginResp, error) {
	//
	err := l.sendHello(in)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RecListenerWhenUserLoginLogic sendHello err:%+v", err)
	}
	return l.getPopListener(in)
}

// 获取登陆弹框的XXX
func (l *RecListenerWhenUserLoginLogic) getPopListener(in *pb.RecListenerWhenUserLoginReq) (*pb.RecListenerWhenUserLoginResp, error) {
	resp := pb.RecListenerWhenUserLoginResp{}
	if in.IsNewUser != db.Enable && in.OrderCnt > 0 { // 新用户或者未下过订单的用户
		return &resp, nil
	}
	var listenerUids []int64
	var err error
	listenerUids, err = l.svcCtx.ListenerRedis.GetRecommendListenerAndIncrScore(l.ctx, rediskey.RedisKeyListenerRecommendWhenUserLogin, 1)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic ListenerRedis.GetRecommendListenerAndIncrScore err:%+v", err)
		return nil, err
	}
	if len(listenerUids) >= 1 {
		var pf *listenerPgModel.ListenerProfile
		pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, listenerUids[0])
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetNewUserRecommendListenerLogic FindOne err:%+v listenerUid:%d", err, listenerUids[0])
			return nil, err
		}
		if tool.IsStringArrayExist(in.AuthKey, listenerkey.TestUserAuthKey) && pf.ListenerUid == cservice.DefaultListenerUid {
			pf.VoiceChatPrice = 100
			pf.TextChatPrice = 100
		}
		resp.RecListener = addRecListenerData(pf)
		return &resp, nil
	}
	return &resp, nil
}

// 发送问候消息
func (l *RecListenerWhenUserLoginLogic) sendHello(in *pb.RecListenerWhenUserLoginReq) (err error) {
	// 超过配置天数 不再发送
	if in.RegDays*24*60*60 > rediskey.RedisKeyUserLoginHelloMsgExpire {
		return nil
	}

	// 判断是否可以发送
	if in.IsNewUser != db.Enable {
		var b bool
		b, err = l.svcCtx.ImRedis.IsHaveNotifyLimit(l.ctx, fmt.Sprintf("%d-%d", in.Uid, notify.TextMsgTypeHelloMsg2), notify.TextMsgTypeHelloMsg2)
		if err != nil {
			return err
		}
		if b {
			return nil
		}
	}

	kqmsg := kqueue.SendHelloWhenUserLoginMessage{
		Uid:       in.Uid,
		IsNewUser: in.IsNewUser,
	}

	var buf []byte
	buf, err = json.Marshal(&kqmsg)
	if err != nil {
		return
	}

	return l.svcCtx.KqSendHelloWhenUserLoginClient.Push(string(buf))
}
