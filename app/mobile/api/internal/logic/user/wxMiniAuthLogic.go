package user

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"
	"net/url"
	"strings"
	"unicode/utf8"
)

type WxMiniAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxMiniAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxMiniAuthLogic {
	return &WxMiniAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxMiniAuthLogic) WxMiniAuth(req *types.WXMiniAuthReq) (resp *types.WXMiniAuthResp, err error) {
	// 校验请求的合法性
	// 061w5RFa1BReYD0633Ha1L0WR60w5RFl
	// 071Wgy000yUrEO1mSD200EYFVu3Wgy0B
	if req.Code == "" || req.IV != "" || req.EncryptedData != "" {
		logx.WithContext(l.ctx).Errorf("WxMiniAuthLogic 参数空校验失败 req:%+v", req)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数校验失败")
	}
	col := utf8.RuneCountInString(req.Code)
	if col != 32 {
		logx.WithContext(l.ctx).Errorf("WxMiniAuthLogic Code长度校验失败 len:%d req:%+v", col, req)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数校验失败")
	}
	if (req.AppVer <= 0 || req.AppVer < l.svcCtx.Config.AppVerConf.MinAppVer) && req.AppVer != -l.svcCtx.Config.AppVerConf.MinAppVer {
		logx.WithContext(l.ctx).Errorf("WxMiniAuthLogic 版本验失败 min ver:%d req:%+v", l.svcCtx.Config.AppVerConf.MinAppVer, req)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数校验失败")
	}

	authResult, err := l.svcCtx.WxMini.GetAuth().Code2SessionContext(l.ctx, req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.WXminiAuthFail, fmt.Sprintf("微信授权失败 %+v", err))
	}

	// 加分布式锁
	//rkey := fmt.Sprintf(rediskey.RedisLockLoginUser, authResult.OpenID)
	//rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	//rl.SetExpire(5)
	//b, err := rl.AcquireCtx(l.ctx)
	//if err != nil {
	//	return nil, err
	//}
	//if !b {
	//	return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	//}
	//defer func() {
	//	_, err2 := rl.ReleaseCtx(l.ctx)
	//	if err2 != nil {
	//		logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
	//		return
	//	}
	//}()

	//
	var channel string
	var cb string
	channel, cb = parseChannel(l.ctx, req)
	in := usercenter.LoginReq{
		AuthKey:   authResult.OpenID,
		AuthType:  userkey.UserAuthTypeWXMini,
		WxUnionId: authResult.UnionID,
		Channel:   channel,
		Cb:        cb,
	}
	//2、login.
	loginResp, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.WXMiniAuthResp{User: &types.UserProfile{}}
	_ = copier.Copy(resp, loginResp)

	resp.LatestAppVer = l.svcCtx.Config.AppVerConf.LatestAppVer
	resp.MinAppVer = l.svcCtx.Config.AppVerConf.MinAppVer
	if req.AppVer <= l.svcCtx.Config.AppVerConf.StoryTabMaxVer {
		resp.StoryTabSwitch = db.Enable
	} else {
		resp.StoryTabSwitch = db.Disable
	}

	return resp, nil
}

func parseChannel(ctx context.Context, req *types.WXMiniAuthReq) (channel string, cb string) {
	if req.Query == "" {
		return
	}
	ul, err := url.QueryUnescape(req.Query)
	if err != nil {
		logx.WithContext(ctx).Errorf("WxMiniAuth parseChannel QueryUnescape err:%+v", err)
		return userkey.GetUserChannelDefault, req.Query
	}
	ul = strings.Replace(ul, "\\u0026", "&", -1)
	//
	//zhihu
	if strings.Contains(ul, "source=zhihu") {
		r := strings.Split(ul, "cb=")
		if len(r) == 2 {
			f := `&source=zhihu&`
			s := strings.Replace(r[1], f, "", -1)
			return userkey.GetUserChannelZhihu, s
		}
	}
	// 其他
	r := strings.Split(ul, "&")
	m := make(map[string]string)
	for idx := 0; idx < len(r); idx++ {
		i := strings.IndexAny(r[idx], "=")
		if i >= 0 {
			m[r[idx][:i]] = r[idx][i+1:]
		}
	}
	return m[userkey.UserChannelKey], ul
}
