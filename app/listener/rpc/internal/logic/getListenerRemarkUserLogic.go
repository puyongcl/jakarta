package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRemarkUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerRemarkUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRemarkUserLogic {
	return &GetListenerRemarkUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX备注的用户
func (l *GetListenerRemarkUserLogic) GetListenerRemarkUser(in *pb.GetListenerRemarkUserReq) (*pb.GetListenerRemarkUserResp, error) {
	rs, err := l.svcCtx.ListenerRemarkUserModel.FindOne(l.ctx, fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid))
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	resp := &pb.GetListenerRemarkUserResp{}
	if rs != nil {
		_ = copier.Copy(resp, rs)
	}
	return resp, nil
}
