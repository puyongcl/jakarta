package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	pbListener "jakarta/app/listener/rpc/pb"
)

type GenListenerContractLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenListenerContractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenListenerContractLogic {
	return &GenListenerContractLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenListenerContractLogic) GenListenerContract(req *types.GenListenerContractReq) (resp *types.GenListenerContractResp, err error) {
	var in pbListener.GenListenerContractReq
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.ListenerRpc.GenListenerContract(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GenListenerContractResp{}
	_ = copier.Copy(resp, rsp)
	return
}
