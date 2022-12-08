package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListListenerProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListListenerProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListListenerProfileLogic {
	return &ListListenerProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListListenerProfileLogic) ListListenerProfile(req *types.GetListenerProfileListReq) (resp *types.GetListenerProfileListResp, err error) {
	var r pb.GetListenerProfileListReq
	_ = copier.Copy(&r, req)
	rsp, err := l.svcCtx.ListenerRpc.AdminGetListenerProfileList(l.ctx, &r)
	resp = &types.GetListenerProfileListResp{List: make([]*types.CheckListenerProfile, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
