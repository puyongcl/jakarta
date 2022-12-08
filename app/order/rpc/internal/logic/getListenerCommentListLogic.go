package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
)

type GetListenerCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerCommentListLogic {
	return &GetListenerCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX评价列表
func (l *GetListenerCommentListLogic) GetListenerCommentList(in *pb.GetListenerCommentListReq) (*pb.GetListenerCommentListResp, error) {
	rs, err := l.svcCtx.ChatOrderModel.FindListenerOrderOpinionList(l.ctx, in.ListenerUid, in.Star, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetListenerCommentListResp{List: make([]*pb.ListenerOrderOpinion, 0)}

	for idx := 0; idx < len(rs); idx++ {
		var val pb.ListenerOrderOpinion
		_ = copier.Copy(&val, rs[idx])
		if rs[idx].CommentTime.Valid {
			val.CommentTime = rs[idx].CommentTime.Time.Format(db.DateFormat)
		}
		//空评价
		if rs[idx].Comment == "" {
			val.Comment = listenerkey.DefaultComment
		}
		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
