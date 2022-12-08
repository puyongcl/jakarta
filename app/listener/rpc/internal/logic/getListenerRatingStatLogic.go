package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRatingStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerRatingStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRatingStatLogic {
	return &GetListenerRatingStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX评价统计情况
func (l *GetListenerRatingStatLogic) GetListenerRatingStat(in *pb.GetListenerRatingStatReq) (*pb.GetListenerRatingStatResp, error) {
	rs, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	//

	resp := &pb.GetListenerRatingStatResp{}
	_ = copier.Copy(resp, rs)

	var rs2 []redis.Pair
	rs2, err = l.svcCtx.ListenerRedis.GetTopCommentTagStat(l.ctx, rediskey.RedisKeyListenerCommentTagStat, in.ListenerUid)
	if err != nil {
		return nil, err
	}

	resp.CommentTagStat = make([]*pb.CommentTagPair, 0)
	var tag int64
	for idx := 0; idx < listenerkey.ShowTopCommentTagCnt; idx++ {
		va := pb.CommentTagPair{}
		va.Tag = tag
		va.Cnt = rs2[idx].Score
		resp.CommentTagStat = append(resp.CommentTagStat, &va)
	}

	return resp, nil
}
