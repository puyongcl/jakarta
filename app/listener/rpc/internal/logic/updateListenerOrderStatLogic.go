package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/rediskey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerOrderStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerOrderStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerOrderStatLogic {
	return &UpdateListenerOrderStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX的统计数据
func (l *UpdateListenerOrderStatLogic) UpdateListenerOrderStat(in *pb.UpdateListenerOrderStatReq) (*pb.UpdateListenerOrderStatResp, error) {
	addData := new(listenerPgModel.AddListenerStat)
	_ = copier.Copy(addData, in)
	err := l.svcCtx.ListenerProfileModel.UpdateListenerStat(l.ctx, addData)
	if err != nil {
		return nil, err
	}
	// 更新XXX评论标签统计
	if len(in.CommentTag) > 0 {
		err = l.svcCtx.ListenerRedis.AddCommentTag(l.ctx, rediskey.RedisKeyListenerCommentTagStat, in.ListenerUid, in.CommentTag)
		if err != nil {
			return nil, err
		}
	}

	// 更新今日接单量
	err = l.svcCtx.StatRedis.AddTodayListenerOrderCnt(l.ctx, in.ListenerUid, in.AddPaidOrderCnt)
	if err != nil {
		return nil, err
	}
	// 更新今日接单金额
	err = l.svcCtx.StatRedis.AddTodayListenerOrderAmount(l.ctx, in.ListenerUid, in.UserPaidOrderAmount)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateListenerOrderStatResp{}, nil
}
