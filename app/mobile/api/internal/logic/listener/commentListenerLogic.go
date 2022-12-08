package listener

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListenerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListenerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListenerLogic {
	return &CommentListenerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListenerLogic) CommentListener(req *types.CommentOrderReq) (resp *types.CommentOrderResp, err error) {
	// 加分布式锁
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid != req.Uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.Uid))
	}
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

	// 评价
	var action int64
	if req.Star == listenerkey.Rating5Star {
		action = orderkey.ChatOrderState5StartRatingAndFinish14
	} else {
		action = orderkey.ChatOrderStateNot5StarWaitConfirm15
	}
	in := pbOrder.DoChatOrderActionReq{
		OrderId:     req.OrderId,
		OperatorUid: req.Uid,
		Comment:     req.Comment,
		Tag:         req.CommentTag,
		Star:        req.Star,
		Action:      action,
		SendMsg:     req.SendMsg,
	}
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	return
}
