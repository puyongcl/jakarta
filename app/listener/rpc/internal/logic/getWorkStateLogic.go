package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkStateLogic {
	return &GetWorkStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX状态设置
func (l *GetWorkStateLogic) GetWorkState(in *pb.GetWorkStateReq) (*pb.GetWorkStateResp, error) {
	rs, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	var val pb.GetWorkStateResp
	if rs != nil {
		_ = copier.Copy(&val, rs)
	}
	return &val, nil
}
