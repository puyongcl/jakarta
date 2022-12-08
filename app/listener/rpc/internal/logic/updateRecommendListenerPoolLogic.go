package logic

import (
	"context"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRecommendListenerPoolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRecommendListenerPoolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRecommendListenerPoolLogic {
	return &UpdateRecommendListenerPoolLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新推荐XXX列表
func (l *UpdateRecommendListenerPoolLogic) UpdateRecommendListenerPool(in *pb.UpdateRecommendListenerPoolReq) (*pb.UpdateRecommendListenerPoolResp, error) {
	// 每天0点后新建一个列表
	rsp, err := l.svcCtx.ListenerProfileModel.FindRecentActive(l.ctx, int(in.RecentDay), uint64(in.Size), []int64{listenerkey.ListenerWorkStateWorking})
	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		var rsp2 []int64
		rsp2, err = l.svcCtx.ListenerProfileModel.FindRecentActive(l.ctx, int(in.RecentDay), uint64(in.Size), []int64{listenerkey.ListenerWorkStateRestingAuto, listenerkey.ListenerWorkStateRestingManual})
		if err != nil {
			return nil, err
		}
		rsp = append(rsp, rsp2...)
	}
	err = l.svcCtx.ListenerRedis.InitRecommendListenerPool(l.ctx, rediskey.RedisKeyListenerRecommendWhenUserLogin, rsp)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ListenerRedis.InitRecommendListenerPool(l.ctx, rediskey.RedisKeyListenerRecommendReplyUserStory, rsp)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ListenerRedis.InitRecommendListenerPool(l.ctx, rediskey.RedisKeyListenerRecommendSendHelloMsgWhenUserLogin, rsp)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateRecommendListenerPoolResp{Cnt: int64(len(rsp))}, nil
}
