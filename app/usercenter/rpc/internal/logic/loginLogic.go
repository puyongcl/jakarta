package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/key/userkey"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrGenerateTokenError = xerr.NewGrpcErrCodeMsg(xerr.TokenGenerateError, "生成token失败")
var ErrUsernamePwdError = xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "账号或密码不正确")

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	userAuth, err := l.svcCtx.UserAuthModel.FindOneByAuthKeyAuthType2(l.ctx, in.AuthKey, in.AuthType)
	if err != nil && err != userPgModel.ErrNotFound {
		return &pb.LoginResp{}, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("根据手机号查询用户信息失败，mobile:%s,err:%+v", in.AuthKey, err))
	}
	if userAuth == nil { // 未注册用户
		if in.UserType == 0 {
			in.UserType = userkey.UserTypeNormalUser
		}
		var regRsp *pb.RegisterResp
		registerLogic := NewRegisterLogic(l.ctx, l.svcCtx)
		in2 := pb.RegisterReq{
			Password:  in.Password,
			AuthKey:   in.AuthKey,
			AuthType:  in.AuthType,
			UserType:  in.UserType,
			Uid:       in.Uid,
			NickName:  in.NickName,
			Avatar:    in.Avatar,
			WxUnionId: in.WxUnionId,
			Channel:   in.Channel,
			Cb:        in.Cb,
		}
		regRsp, err = registerLogic.Register(&in2)
		if err != nil {
			return &pb.LoginResp{}, err
		}

		resp := pb.LoginResp{
			AccessToken:  regRsp.AccessToken,
			AccessExpire: regRsp.AccessExpire,
			RefreshAfter: regRsp.RefreshAfter,
			User:         regRsp.User,
			OpenId:       in.AuthKey,
			AccountState: regRsp.AccountState,
			UserType:     regRsp.UserType,
			UserSign:     regRsp.UserSign,
			IsNewUser:    db.Enable,
			Channel:      in.Channel,
		}

		// 是否关注服务号
		resp.IsFollowWxFwh = IsFollowWxFwh(l.ctx, l.svcCtx, in.WxUnionId)

		return &resp, nil
	}

	if !tool.IsInt64ArrayExist(userAuth.UserType, userkey.AllowLoginUserType) {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.UserTypeNotAllowLogin, "不能登陆的账号")
	}

	return doLogin(l.ctx, l.svcCtx, in, userAuth)
}

func IsFollowWxFwh(ctx context.Context, svcCtx *svc.ServiceContext, unionId string) (is int64) {
	if unionId == "" {
		return
	}
	data, err := svcCtx.UserWechatInfoModel.FindOne(ctx, unionId)
	if err != nil && err != userPgModel.ErrNotFound {
		logx.WithContext(ctx).Errorf("LoginLogic IsFollowWxFwh err:%+v", err)
		return
	}

	if data == nil {
		return
	}
	return data.FwhState
}

func doLogin(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.LoginReq, userAuth *userPgModel.UserAuth) (*pb.LoginResp, error) {
	//
	uid := userAuth.Uid
	var err error
	switch in.AuthType {
	case userkey.UserAuthTypePasswd:
		if !(tool.Md5ByString(in.Password) == userAuth.Password) {
			return &pb.LoginResp{}, ErrUsernamePwdError
		}
	case userkey.UserAuthTypeWXMini: // no need check
		break
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "unsupported type")
	}

	// 检查封禁状态
	if userAuth.FreeTime.Valid {
		if userAuth.FreeTime.Time.After(time.Now()) {
			userAuth.AccountState = userkey.AccountStateBan
		}
	}

	// Generate the token, so that the service doesn't call rpc internally
	tokenResp := &pb.GenerateTokenResp{
		AccessToken:  "",
		AccessExpire: 0,
		RefreshAfter: 0,
	}
	if userAuth.AccountState == userkey.AccountStateNormal {
		generateTokenLogic := NewGenerateTokenLogic(ctx, svcCtx)
		tokenResp, err = generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
			Uid:         uid,
			UserChannel: userAuth.Channel,
			AppVer:      101010,
			AuthType:    userAuth.AuthType,
			AuthKey:     userAuth.AuthKey,
		})
		if err != nil {
			return nil, err
		}
	}

	// get userAuth profile
	getUserInfoLogic := NewGetUserInfoLogic(ctx, svcCtx)
	userInfoResp, err := getUserInfoLogic.GetUserInfo(&pb.GetUserProfileReq{Uid: uid})
	if err != nil {
		return nil, err
	}

	// get tim user signature
	var sig string
	sig, err = svcCtx.UserRedis.GetTimUserSignature(ctx, uid)
	if err != nil {
		return nil, err
	}
	if sig == "" && userAuth.AccountState == userkey.AccountStateNormal {
		registerLogic := NewRegisterLogic(ctx, svcCtx)
		sig, err = registerLogic.GetTimUserSignature(uid)
		if err != nil {
			return nil, err
		}
	}
	if userAuth.AccountState != userkey.AccountStateNormal {
		sig = ""
	}

	resp := usercenter.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
		User:         userInfoResp.User,
		UserSign:     sig,
		AccountState: userAuth.AccountState,
		UserType:     userAuth.UserType,
		OpenId:       in.AuthKey,
		IsNewUser:    db.Disable,
		Channel:      userAuth.Channel,
	}
	if userAuth.AccountState == userkey.AccountStateBan {
		resp.FreeTime = userAuth.FreeTime.Time.Format(db.DateTimeFormat)
		resp.BanReason = userAuth.BanReason
	}
	// 是否关注服务号
	resp.IsFollowWxFwh = IsFollowWxFwh(ctx, svcCtx, in.WxUnionId)

	return &resp, nil
}
