package listener

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerFeedbackUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListenerFeedbackUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerFeedbackUserLogic {
	return &ListenerFeedbackUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListenerFeedbackUserLogic) ListenerFeedbackUser(req *types.FeedbackOrderReq) (resp *types.FeedbackOrderResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	if req.OrderId == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	if req.Feedback == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.OrderErrorFeedbackNotAllowEmpty, "没有填写反馈内容")
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

	_, err = l.svcCtx.OrderRpc.FeedbackOrder(l.ctx, &pbOrder.FeedbackOrderReq{
		OrderId:     req.OrderId,
		Uid:         req.Uid,
		ListenerUid: req.ListenerUid,
		Feedback:    req.Feedback,
		SendMsg:     req.SendMsg,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.FeedbackOrderResp{}
	return
}
