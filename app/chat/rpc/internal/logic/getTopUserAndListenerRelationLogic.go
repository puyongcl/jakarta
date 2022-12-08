package logic

import (
	"context"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopUserAndListenerRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTopUserAndListenerRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopUserAndListenerRelationLogic {
	return &GetTopUserAndListenerRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取交互最频繁的几位XXX
func (l *GetTopUserAndListenerRelationLogic) GetTopUserAndListenerRelation(in *pb.GetTopUserAndListenerRelationReq) (*pb.GetTopUserAndListenerRelationResp, error) {
	rs, err := l.svcCtx.UserListenerRelationModel.FindTopScoreList(l.ctx, in.Uid, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}

	resp := pb.GetTopUserAndListenerRelationResp{}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.UserAndListenerRelation
		val.ListenerUid = rs[idx].ListenerUid
		val.Score = rs[idx].TotalScore

		resp.List = append(resp.List, &val)
	}
	return &resp, nil
}
