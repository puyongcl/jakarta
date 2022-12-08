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

type AddStoryReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStoryReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStoryReplyLogic {
	return &AddStoryReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStoryReplyLogic) AddStoryReply(req *types.AddStoryReplyReq) (resp *types.AddStoryReplyResp, err error) {
	if req.ReplyText == "" && req.ReplyVoice == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "回复不能为空")
	}
	if req.ReplyText != "" && utf8.RuneCountInString(req.ReplyText) < 20 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "回复不能少于20字")
	}
	if utf8.RuneCountInString(req.ReplyText) > 600 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.BbsErrorStoryVerifyError, "回复不能超过600字")
	}

	// 加分布式锁
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid != req.ListenerUid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
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

	var in pb.AddStoryReplyReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.BbsRpc.AddStoryReply(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.AddStoryReplyResp{Id: rsp.Id}
	return
}
