package logic

import (
	"context"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerOnlineStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerOnlineStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerOnlineStateLogic {
	return &UpdateListenerOnlineStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX的登陆状态
func (l *UpdateListenerOnlineStateLogic) UpdateListenerOnlineState(in *pb.UpdateListenerOnlineStateReq) (*pb.UpdateListenerOnlineStateResp, error) {
	err := l.svcCtx.ListenerProfileModel.UpdateOnlineState(l.ctx, in.ListenerUid, in.State)
	if err != nil {
		return nil, err
	}

	if in.TodayLoginCnt == 1 { // 第一次登陆时加入 新用户推荐列表
		// 初始化今日排名
		ulo := NewUpdateListenerOrderStatLogic(l.ctx, l.svcCtx)
		_, err = ulo.UpdateListenerOrderStat(&pb.UpdateListenerOrderStatReq{ListenerUid: in.ListenerUid})
		if err != nil {
			return nil, err
		}

		// 加入新用户推荐
		var pf *listenerPgModel.ListenerProfile
		pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil {
			return nil, err
		}

		if pf.WorkState == listenerkey.ListenerWorkStateWorking { // 不是工作中 不加入
			err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendWhenUserLogin, in.ListenerUid)
			if err != nil {
				return nil, err
			}
			err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendReplyUserStory, in.ListenerUid)
			if err != nil {
				return nil, err
			}
			err = l.svcCtx.ListenerRedis.ADDNewUserRecommendListenerOne(l.ctx, rediskey.RedisKeyListenerRecommendSendHelloMsgWhenUserLogin, in.ListenerUid)
			if err != nil {
				return nil, err
			}
		}

	}

	return &pb.UpdateListenerOnlineStateResp{}, nil
}
