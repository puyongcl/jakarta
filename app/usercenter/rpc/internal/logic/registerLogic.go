package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/key/tencentcloudkey"
	tim2 "jakarta/common/key/timkey"
	"jakarta/common/key/userkey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	tim3 "jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	var err error
	// gen uid
	var uid int64
	if in.Uid != 0 {
		uid = in.Uid
	} else {
		uid, err = l.svcCtx.UserRedis.IncrUidIdx(l.ctx)
		if err != nil {
			return nil, err
		}
		uid = uid + userkey.UidStart
	}

	// gen name
	var name string
	if in.NickName != "" {
		name = in.NickName
	} else {
		name = userkey.GetDefaultUsername(uid)
	}

	// gen avatar
	var avatar string
	if in.Avatar != "" {
		avatar = in.Avatar
	} else {
		avatar = userkey.GetDefaultAvatar(uid)
	}

	if in.Channel == "" {
		in.Channel = userkey.GetUserChannelDefault
	}

	err = l.svcCtx.UserAuthModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		user := new(userPgModel.UserProfile)
		user.Uid = uid
		user.Nickname = name
		user.Avatar = avatar
		user.Birthday = sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}

		_, err = l.svcCtx.UserProfileModel.InsertTrans(ctx, session, user)
		if err != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user Insert err:%+v,uid:%d", err, uid))
		}

		//
		userAuth := new(userPgModel.UserAuth)
		userAuth.Uid = uid
		userAuth.AuthKey = in.AuthKey
		userAuth.AuthType = in.AuthType
		userAuth.AccountState = userkey.AccountStateNormal
		userAuth.UserType = in.UserType
		userAuth.Channel = in.Channel
		if len(in.Password) > 0 && in.AuthType == userkey.UserAuthTypePasswd {
			userAuth.Password = tool.Md5ByString(in.Password)
		}

		_, err = l.svcCtx.UserAuthModel.InsertTrans(ctx, session, userAuth)
		if err != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user_auth Insert err:%+v,uid:%d", err, uid))
		}

		//
		userLoginState := new(userPgModel.UserLoginState)
		userLoginState.LoginState = userkey.Unknown // 默认未知登陆状态 存在已注册 未登陆过im的情况
		userLoginState.LoginCntSum = 1
		userLoginState.LoginCntToday = 1
		userLoginState.Uid = uid
		userLoginState.LoginTime = time.Now()
		userLoginState.OfflineTime = time.Now()
		_, err = l.svcCtx.UserLoginStateModel.InsertTrans(ctx, session, userLoginState)
		if err != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user_auth Insert err:%+v,uid:%d", err, uid))
		}

		//
		userStat := new(userPgModel.UserStat)
		userStat.Uid = uid
		userStat.NoCondRefundCnt = userkey.NoCondRefundCnt
		_, err = l.svcCtx.UserStatModel.InsertTrans(ctx, session, userStat)
		if err != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user_stat Insert err:%+v,uid:%d", err, uid))
		}
		// 微信平台信息
		if in.WxUnionId != "" {
			uwi := new(userPgModel.UserWechatInfo)
			*uwi = userPgModel.UserWechatInfo{
				Uid:       uid,
				FwhOpenid: "",
				UnionId:   in.WxUnionId,
				FwhState:  0,
			}
			if in.AuthType == userkey.UserAuthTypeWXMini {
				uwi.MpOpenid = in.AuthKey
			}
			_, err = l.svcCtx.UserWechatInfoModel.InsertOrUpdateMPTrans(ctx, session, uwi)
			if err != nil {
				return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user_wechat_info Insert err:%+v,uid:%d", err, uid))
			}
		}

		// 获客渠道事件回传
		if in.Cb != "" {
			ucb := new(userPgModel.UserChannelCallback)
			*ucb = userPgModel.UserChannelCallback{
				Uid:     uid,
				Channel: in.Channel,
				Cb:      in.Cb,
			}
			_, err = l.svcCtx.UserChannelCallbackModel.InsertTrans(ctx, session, ucb)
			if err != nil {
				return xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("Register db user_channel_callback InsertTrans err:%+v,uid:%d", err, uid))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	//2、Generate the token, so that the service doesn't call rpc internally
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		Uid:         uid,
		UserChannel: in.Channel,
		AppVer:      101010,
		AuthKey:     in.AuthKey,
		AuthType:    in.AuthType,
	})
	if err != nil {
		return nil, ErrGenerateTokenError
	}

	// get tim user sig
	sig, err := l.GetTimUserSignature(uid)
	if err != nil {
		return nil, err
	}

	// 注册IM并延迟发送新用户im消息
	l.regIm(uid, name, avatar, in.UserType)

	// 获客渠道新用户注册上报
	l.regEvent(uid, in.Channel, in.Cb)

	return &usercenter.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
		User: &pb.UserProfile{
			CreateTime: time.Now().Format(db.DateTimeFormat),
			Uid:        uid,
			Nickname:   name,
			Avatar:     avatar,
		},
		UserSign:     sig,
		AccountState: userkey.AccountStateNormal,
		UserType:     in.UserType,
	}, nil
}

// 上报获客渠道新注册用户
func (l *RegisterLogic) regEvent(uid int64, channel string, cb string) {
	if cb == "" {
		return
	}
	//
	var kqMsg *kqueue.UploadUserEventMessage
	switch channel {
	case userkey.GetUserChannelZhihu:
		kqMsg = &kqueue.UploadUserEventMessage{
			Uid:   uid,
			Cb:    cb,
			Event: userkey.ZhihuUploadEventReg,
			Value: "",
			Stamp: fmt.Sprintf("%d", time.Now().Unix()),
		}
	default:
		return
	}

	buf, err := json.Marshal(kqMsg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RegisterLogic regEvent json Marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueUploadUserEventClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RegisterLogic regEvent push err:%+v", err)
		return
	}
	return
}

func (l *RegisterLogic) regIm(uid int64, nick, avatar string, userType int64) {
	// 注册IM账号
	err := l.svcCtx.TimClient.AccountImport(uid, tencentcloudkey.CDNBasePath+avatar, nick)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RegisterLogic AccountImport uid:%d error:%+v", uid, err)
		return
	}
	if userType != userkey.UserTypeNormalUser {
		return
	}
	// 发送新用户消息
	msg := kqueue.SendImDefineMessage{
		FromUid: notify.TimSystemNotifyUid,
		ToUid:   uid,
		MsgType: notify.DefineNotifyMsgTypeSystemMsg18,
		Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle18,
		Text:    notify.DefineNotifyMsgTemplateSystemMsg18,
		Val1:    "",
		Val2:    "",
		Sync:    tim3.TimMsgSyncFromNo,
	}
	var buf []byte
	buf, err = json.Marshal(msg)
	if err != nil {
		return
	}
	var payload []byte
	payload, err = json.Marshal(jobtype.DeferSendImMsgPayload{KqMsgBuf: buf})
	if err != nil {
		return
	}
	_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferSendImMsg, payload), asynq.ProcessIn(notify.DeferSendImMsgSecond*time.Second))
	if err != nil {
		return
	}
}

func (l *RegisterLogic) GetTimUserSignature(uid int64) (sig string, err error) {
	// gen tim user signature
	sig, err = tim3.GenUserSig(uid, l.svcCtx.Config.TimConf.SDKAPPID, l.svcCtx.Config.TimConf.IMKEY, tim2.TimUserSigExpire)
	if err != nil {
		return "", xerr.NewGrpcErrCodeMsg(xerr.GetUserSignFail, fmt.Sprintf("user:%d get im sign err : %+v", uid, err))
	}

	// set user signature in redis
	err = l.svcCtx.UserRedis.SetTimUserSignature(l.ctx, uid, sig, 7*24*60*60)
	if err != nil {
		return "", err
	}
	return
}
