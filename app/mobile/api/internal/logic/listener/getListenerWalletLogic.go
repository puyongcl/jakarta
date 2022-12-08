package listener

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerWalletLogic {
	return &GetListenerWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerWalletLogic) GetListenerWallet(req *types.GetListenerWalletReq) (resp *types.GetListenerWalletResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	var in pbListener.GetListenerWalletReq
	in.ListenerUid = req.ListenerUid
	rs, err := l.svcCtx.ListenerRpc.GetListenerWallet(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerWalletResp{}
	_ = copier.Copy(resp, rs)
	return
}
