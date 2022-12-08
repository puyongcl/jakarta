package bbs

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/bbs/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStoryReplyLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStoryReplyLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStoryReplyLikeLogic {
	return &AddStoryReplyLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStoryReplyLikeLogic) AddStoryReplyLike(req *types.AddLikeStoryReplyReq) (resp *types.AddLikeStoryReplyResp, err error) {
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

	var in pb.AddLikeStoryReplyReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.AddLikeStoryReply(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.AddLikeStoryReplyResp{LikeCnt: rsp.LikeCnt}

	return
}
