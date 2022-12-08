package listener

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"
	"unicode/utf8"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditListenerWordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditListenerWordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditListenerWordsLogic {
	return &EditListenerWordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditListenerWordsLogic) EditListenerWords(req *types.EditListenerWordsReq) (resp *types.EditListenerWordsResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}

	if utf8.RuneCountInString(req.Words1) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words2) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words3) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words4) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words5) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words6) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words7) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words8) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words9) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
	}
	if utf8.RuneCountInString(req.Words10) > listenerkey.MaxListenerWordsCount {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorEditWords, fmt.Sprintf("字数超过%d", listenerkey.MaxListenerWordsCount))
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

	var in pb.EditListenerWordsReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.ListenerRpc.EditListenerWords(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.EditListenerWordsResp{}
	return
}
