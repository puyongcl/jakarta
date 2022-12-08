package listener

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommitNewListenerCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommitNewListenerCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommitNewListenerCheckLogic {
	return &CommitNewListenerCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommitNewListenerCheckLogic) CommitNewListenerCheck(req *types.NewListenerCommitCheckReq) (resp *types.NewListenerCommitCheckResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.Uid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.Uid))
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

	var in pb.CommitCheckNewListenerProfileReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.ListenerRpc.CommitCheckNewListenerProfile(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.NewListenerCommitCheckResp{}
	return
}
