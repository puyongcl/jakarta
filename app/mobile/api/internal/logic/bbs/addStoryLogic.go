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
	"unicode/utf8"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddStoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStoryLogic {
	return &AddStoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStoryLogic) AddStory(req *types.AddStoryReq) (resp *types.AddStoryResp, err error) {
	if req.Tittle == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "标题不能为空")
	}
	if req.Content == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "内容不能为空")
	}
	if utf8.RuneCountInString(req.Content) <= 20 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "内容不能少于20字")
	}
	if utf8.RuneCountInString(req.Tittle) > 60 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "标题不能超过60字")
	}

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

	var in pb.AddStoryReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.AddStory(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.AddStoryResp{Id: rsp.Id}
	return
}
