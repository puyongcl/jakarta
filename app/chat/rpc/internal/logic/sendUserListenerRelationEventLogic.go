package logic

import (
	"context"
	"fmt"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"

	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendUserListenerRelationEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendUserListenerRelationEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendUserListenerRelationEventLogic {
	return &SendUserListenerRelationEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  用户和XXX交互事件
func (l *SendUserListenerRelationEventLogic) SendUserListenerRelationEvent(in *pb.SendUserListenerRelationEventReq) (*pb.SendUserListenerRelationEventResp, error) {
	add := chatkey.GetAddScore(in.EventType)
	if add <= 0 {
		return nil, nil
	}
	id := fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid)
	err := l.svcCtx.UserListenerRelationModel.AddScore(l.ctx, add, id)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("SendUserListenerRelationEventLogic AddScore in:%+v err:%+v", in, err)
	}
	return &pb.SendUserListenerRelationEventResp{}, nil
}
