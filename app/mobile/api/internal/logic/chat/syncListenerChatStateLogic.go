package chat

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncListenerChatStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncListenerChatStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncListenerChatStateLogic {
	return &SyncListenerChatStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncListenerChatStateLogic) SyncListenerChatState(req *types.SyncListenerChatStateReq) (resp *types.SyncListenerChatStateResp, err error) {
	resp, err = l.DoSyncListenerChatState(req)
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Infof("SyncListenerChatStateLogic caller uid:%d,action:%d resp:%+v err:%+v", uid, req.Action, resp, err)
	} else {
		logx.WithContext(l.ctx).Infof("SyncListenerChatStateLogic caller uid:%d,action:%d resp:%+v", uid, req.Action, resp)
	}

	return
}

func (l *SyncListenerChatStateLogic) DoSyncListenerChatState(req *types.SyncListenerChatStateReq) (resp *types.SyncListenerChatStateResp, err error) {
	// 非用户账户 直接返回默认数据
	if (req.ListenerUid != 0 && req.ListenerUid < userkey.UidStart) || req.Uid == req.ListenerUid {
		return &types.SyncListenerChatStateResp{
			WorkState:          listenerkey.ListenerWorkStateWorking,
			FreeChatCnt:        10000000,
			TextChatExpiryTime: "",
			VoiceChatMinute:    0,
			ChatState:          0,
			ListenerChatState:  0,
		}, nil
	}

	// 加分布式锁
	if req.Action == chatkey.ChatAction13 {
		uid := ctxdata.GetUidFromCtx(l.ctx)
		rkey := fmt.Sprintf(rediskey.RedisLockUser, uid)
		rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
		rl.SetExpire(2)
		b, err := rl.AcquireCtx(l.ctx)
		if err != nil {
			return nil, err
		}
		if !b {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
		}
		defer func() {
			_, err2 := rl.ReleaseCtx(l.ctx)
			if err2 != nil {
				logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
				return
			}
		}()
	}

	switch req.Action {
	case chatkey.ChatAction1, chatkey.ChatAction2, chatkey.ChatAction12, chatkey.ChatAction13, chatkey.ChatAction14:
		return l.SyncChat(req)
	case chatkey.ChatAction11:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "wrong action")
	case chatkey.ChatAction3, chatkey.ChatAction4, chatkey.ChatAction5, chatkey.ChatAction6, chatkey.ChatAction7, chatkey.ChatAction8, chatkey.ChatAction9, chatkey.ChatAction10:
		return l.SyncVoiceChat(req)

	default:

	}
	return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "wrong action type")
}

func (l *SyncListenerChatStateLogic) SyncChat(req *types.SyncListenerChatStateReq) (resp *types.SyncListenerChatStateResp, err error) {
	var rs1 *pbChat.SyncChatStateResp
	var rs2 *pbListener.GetWorkStateResp
	err = mr.Finish(func() error {
		var in pbChat.SyncChatStateReq
		_ = copier.Copy(&in, req)
		var err2 error
		rs1, err2 = l.svcCtx.ChatRpc.SyncChatState(l.ctx, &in)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var in pbListener.GetWorkStateReq
		in.ListenerUid = req.ListenerUid
		var err2 error
		rs2, err2 = l.svcCtx.ListenerRpc.GetWorkState(l.ctx, &in)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SyncListenerChatStateResp{}
	_ = copier.Copy(resp, rs1)
	resp.WorkState = rs2.WorkState
	resp.OnlineState = rs2.OnlineState
	return
}

func (l *SyncListenerChatStateLogic) SyncVoiceChat(req *types.SyncListenerChatStateReq) (resp *types.SyncListenerChatStateResp, err error) {
	// 根据openid查询用户id
	var u1 *pbUser.UserAuth
	var u2 *pbUser.UserAuth
	if req.OpenId1 != "" || req.OpenId2 != "" {
		err = mr.Finish(func() error {
			if req.OpenId1 == "" {
				return nil
			}
			var in pbUser.GetUserAuthByAuthKeyReq
			in.AuthKey = req.OpenId1
			in.AuthType = userkey.UserAuthTypeWXMini
			var rs *pbUser.GetUserAuthByAuthKeyResp
			var err2 error
			rs, err2 = l.svcCtx.UsercenterRpc.GetUserAuthByAuthKey(l.ctx, &in)
			if err2 != nil {
				return err2
			}
			u1 = rs.UserAuth
			return nil
		}, func() error {
			if req.OpenId2 == "" {
				return nil
			}
			var in pbUser.GetUserAuthByAuthKeyReq
			in.AuthKey = req.OpenId2
			in.AuthType = userkey.UserAuthTypeWXMini
			var rs *pbUser.GetUserAuthByAuthKeyResp
			var err2 error
			rs, err2 = l.svcCtx.UsercenterRpc.GetUserAuthByAuthKey(l.ctx, &in)
			if err2 != nil {
				return err2
			}
			u2 = rs.UserAuth
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	//
	var in2 pbChat.SyncChatStateReq
	in2.Action = req.Action
	if u1 != nil && u1.UserType == userkey.UserTypeNormalUser {
		in2.Uid = u1.Uid
	}
	if u1 != nil && u1.UserType == userkey.UserTypeListener {
		in2.ListenerUid = u1.Uid
	}
	if u2 != nil && u2.UserType == userkey.UserTypeNormalUser {
		in2.Uid = u2.Uid
	}
	if u2 != nil && u2.UserType == userkey.UserTypeListener {
		in2.ListenerUid = u2.Uid
	}
	if in2.Uid == 0 && in2.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "uid not exist")
	}

	//
	var rs1 *pbChat.SyncChatStateResp
	var rs2 *pbListener.GetWorkStateResp

	rs1, err = l.svcCtx.ChatRpc.SyncChatState(l.ctx, &in2)
	if err != nil {
		return
	}

	var in pbListener.GetWorkStateReq
	in.ListenerUid = rs1.ListenerUid
	rs2, err = l.svcCtx.ListenerRpc.GetWorkState(l.ctx, &in)
	if err != nil {
		return
	}

	resp = &types.SyncListenerChatStateResp{}
	_ = copier.Copy(resp, rs1)
	resp.WorkState = rs2.WorkState
	return
}
