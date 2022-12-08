package cacheModel

import (
	"context"
	"jakarta/app/usercenter/rpc/pb"
	"strconv"
)

type UserInfo struct {
	Uid      int64 `json:"uid"`
	UserType int64 `json:"userType"`
}

func (m *IMMemoryCache) SetUserInfo(uid string, info *UserInfo) {
	m.c.Set(uid, info)
}

func (m *IMMemoryCache) GetUserInfo(ctx context.Context, uid int64) (*UserInfo, error) {
	v, err := m.c.Take(strconv.FormatInt(uid, 10), func() (interface{}, error) {
		return m.getUserInfo(ctx, uid)
	})
	if err != nil {
		return nil, err
	}
	v2, ok := v.(*UserInfo)
	if ok {
		return v2, nil
	}
	return m.getUserInfo(ctx, uid)
}

func (m *IMMemoryCache) getUserInfo(ctx context.Context, uid int64) (*UserInfo, error) {
	rsp, err := m.userRpc.GetUserAuthByUserId(ctx, &pb.GetUserAuthByUserIdReq{Uid: uid})
	if err != nil {
		return nil, err
	}
	r := UserInfo{
		Uid:      uid,
		UserType: rsp.UserAuth.UserType,
	}
	return &r, nil
}
