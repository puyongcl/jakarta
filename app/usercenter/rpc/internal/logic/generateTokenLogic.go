package logic

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, in.Uid, in.UserChannel, in.AppVer, in.AuthKey, in.AuthType)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("getJwtToken err userId:%d , err:%v", in.Uid, err))
	}

	return &pb.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64, userChannel string, appVer int64, authKey, authType string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[ctxdata.CtxKeyJwtUserId] = userId
	claims[ctxdata.CtxKeyJwtUserChannel] = userChannel
	claims[ctxdata.CtxKeyJwtAppVer] = appVer
	claims[ctxdata.CtxKeyJwtAuthKey] = authKey
	claims[ctxdata.CtxKeyJwtAuthType] = authType

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
