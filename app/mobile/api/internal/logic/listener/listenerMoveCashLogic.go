package listener

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerMoveCashLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListenerMoveCashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerMoveCashLogic {
	return &ListenerMoveCashLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListenerMoveCashLogic) ListenerMoveCash(req *types.ListenerMoveCashReq) (resp *types.ListenerMoveCashResp, err error) {
	if req.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty req")
	}
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

	in := pbListener.UpdateListenerWalletReq{
		ListenerUid: req.ListenerUid,
		Amount:      req.Amount,
		SettleType:  listenerkey.ListenerSettleTypeApplyCash,
		OutTime:     time.Now().Format(db.DateTimeFormat),
		Remark:      "",
	}
	_, err = l.svcCtx.ListenerRpc.UpdateListenerWallet(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	return
}
