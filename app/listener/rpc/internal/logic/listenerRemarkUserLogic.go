package logic

import (
	"context"
	"fmt"
	"jakarta/common/key/db"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListenerRemarkUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListenerRemarkUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListenerRemarkUserLogic {
	return &ListenerRemarkUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  XXX备注用户
func (l *ListenerRemarkUserLogic) ListenerRemarkUser(in *pb.ListenerRemarkUserReq) (*pb.ListenerRemarkUserResp, error) {
	err := l.svcCtx.ListenerRemarkUserModel.InsertOrUpdateUserRemark(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid), in.Uid, in.ListenerUid, in.Remark, in.UserDesc)
	if err != nil {
		return nil, err
	}
	return &pb.ListenerRemarkUserResp{}, nil
}
