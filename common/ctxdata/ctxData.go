package ctxdata

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

// CtxKeyJwtUserId get uid from ctx
const CtxKeyJwtUserId = "jwtUserId"
const CtxKeyJwtUserChannel = "jwtUserChannel"
const CtxKeyJwtAppVer = "jwtUserAppVer"
const CtxKeyJwtAuthKey = "jwtUserAuthKey"
const CtxKeyJwtAuthType = "jwtUserAuthType"

func GetUidFromCtx(ctx context.Context) int64 {
	var uid int64
	jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number)
	if ok {
		var err error
		uid, err = jsonUid.Int64()
		if err != nil {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}
	return uid
}

func GetUserChannelFromCtx(ctx context.Context) string {
	var userChannel string
	var ok bool
	userChannel, ok = ctx.Value(CtxKeyJwtUserChannel).(string)
	if !ok {
		logx.WithContext(ctx).Errorf("GetUserChannelFromCtx user channel is not ok")
	}
	return userChannel
}

func GetUserAuthKeyFromCtx(ctx context.Context) string {
	var authKey string
	var ok bool
	authKey, ok = ctx.Value(CtxKeyJwtAuthKey).(string)
	if !ok {
		logx.WithContext(ctx).Errorf("GetUserAuthKeyFromCtx auth key is not ok")
	}
	return authKey
}

func GetUserAuthTypeFromCtx(ctx context.Context) string {
	var authType string
	var ok bool
	authType, ok = ctx.Value(CtxKeyJwtAuthType).(string)
	if !ok {
		logx.WithContext(ctx).Errorf("GetUserAuthTypeFromCtx auth type is not ok")
	}
	return authType
}

func GetAppVerFromCtx(ctx context.Context) int64 {
	var ver int64
	jsonV, ok := ctx.Value(CtxKeyJwtAppVer).(json.Number)
	if ok {
		var err error
		ver, err = jsonV.Int64()
		if err != nil {
			logx.WithContext(ctx).Errorf("GetAppVerFromCtx err : %+v", err)
		}
	}
	return ver
}
