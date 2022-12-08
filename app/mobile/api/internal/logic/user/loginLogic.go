package user

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/cservice"
	"jakarta/common/key/userkey"
	"jakarta/common/notify"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strings"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if !tool.IsStringArrayExist(req.AuthType, userkey.AuthType) {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "auth type error")
	}

	// 校验注册客服号 不是客服号注册 则清空
	if !tool.IsInt64ArrayExist(req.UserType, []int64{userkey.UserTypeCustomerService, userkey.UserTypeNotify, userkey.UserTypeAdmin}) {
		req.Uid = 0
		req.Avatar = ""
		req.NickName = ""
		req.UserType = 0
	}
	// 注册特殊账号 需要校验CODE
	if req.UserType != 0 || req.Uid != 0 || req.Avatar != "" || req.NickName != "" {
		if req.Code != "JAKARTA" {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "参数错误")
		}
		// 校验通过
		// 客服号uid不能小于 20 大于 普通用户id起点
		if req.UserType == userkey.UserTypeCustomerService && (req.Uid > userkey.UidCSMax || req.Uid < cservice.CustomerServiceAccountUid) {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "uid不符合规则")
		}
		// 消息通知号uid不能大于20 小于10 起点
		if req.UserType == userkey.UserTypeNotify && (req.Uid < notify.TimSystemNotifyUid || req.Uid >= cservice.CustomerServiceAccountUid) {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "uid不符合规则")
		}
		// 管理员通知号uid不能大于20 小于10 起点
		if req.UserType == userkey.UserTypeAdmin && (req.Uid < userkey.UidCSMax || req.Uid >= userkey.UidAdminMax) {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "uid不符合规则")
		}
		// 校验头像
		if strings.Contains(req.Avatar, "http") {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "头像不符合规则")
		}
		// 校验昵称
		if len(req.NickName) > 24 {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.UserRegInfoError, "昵称不符合规则")
		}
	}

	// 设置特殊账号默认值
	switch req.UserType {
	case userkey.UserTypeCustomerService:
		req.AuthKey = fmt.Sprintf(cservice.DefaultCSAuthKey, req.Uid)
		req.AuthType = userkey.UserAuthTypePasswd
		req.Password = cservice.DefaultCSAuthPassword
		if req.NickName == "" {
			req.NickName = cservice.DefaultNickName
		}
	default:

	}

	loginResp, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{
		AuthType: req.AuthType,
		AuthKey:  req.AuthKey,
		Password: req.Password,
		Uid:      req.Uid,
		NickName: req.NickName,
		Avatar:   req.Avatar,
		UserType: req.UserType,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.LoginResp{User: &types.UserProfile{}}
	_ = copier.Copy(resp, loginResp)
	resp.LatestAppVer = l.svcCtx.Config.AppVerConf.LatestAppVer
	resp.MinAppVer = l.svcCtx.Config.AppVerConf.MinAppVer
	return resp, nil
}
