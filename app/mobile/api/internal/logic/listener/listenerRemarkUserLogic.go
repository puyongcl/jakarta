package listener

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerRemarkUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListenerRemarkUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerRemarkUserLogic {
	return &ListenerRemarkUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListenerRemarkUserLogic) ListenerRemarkUser(req *types.ListenerRemarkUserReq) (resp *types.ListenerRemarkUserResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}

	// 加分布式锁
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

	_, err = l.svcCtx.ListenerRpc.ListenerRemarkUser(l.ctx, &pb.ListenerRemarkUserReq{
		Uid:         req.Uid,
		ListenerUid: req.ListenerUid,
		Remark:      req.Remark,
		UserDesc:    req.UserDesc,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ListenerRemarkUserResp{}
	return
}
