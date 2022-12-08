package contract1021

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignContract1021Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignContract1021Logic(ctx context.Context, svcCtx *svc.ServiceContext) *SignContract1021Logic {
	return &SignContract1021Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignContract1021Logic) SignContract1021(req *types.SignContract1021Req) (resp *types.SignContract1021Resp, err error) {
	if req.ContractId == "" || req.SignName == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	// 加分布式锁
	rkey := fmt.Sprintf(rediskey.RedisLockGenContract1021, req.ContractId)
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

	st, err := l.svcCtx.Contract1021Model.QuerySignTime(l.ctx, req.ContractId)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SignContract1021Logic QuerySignTime req:%+v err:%+v", req, err)
		return nil, err
	}

	if st != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "无效操作")
	}

	err = l.svcCtx.Contract1021Model.Sign(l.ctx, req.ContractId, req.SignName)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SignContract1021Logic Sign req:%+v err:%+v", req, err)
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "系统错误")
	}
	return
}
