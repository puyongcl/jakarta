package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
)

type GetListenerProfileByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerProfileByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerProfileByUserLogic {
	return &GetListenerProfileByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  普通用户获取XXX资料
func (l *GetListenerProfileByUserLogic) GetListenerProfileByUser(in *pb.GetListenerProfileByUserReq) (*pb.GetListenerProfileByUserResp, error) {
	rs, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	var pf pb.UserSeeListenerProfile
	pf.Specialties = make([]int64, 0)
	if rs != nil {
		_ = copier.Copy(&pf, rs)
	}

	return &pb.GetListenerProfileByUserResp{Profile: &pf}, nil
}
