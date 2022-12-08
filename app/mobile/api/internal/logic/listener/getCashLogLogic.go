package listener

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCashLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCashLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCashLogLogic {
	return &GetCashLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCashLogLogic) GetCashLog(req *types.GetListenerCashListReq) (resp *types.GetListenerCashListResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	var in pbListener.GetListenerCashLogReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.ListenerRpc.GetListenerCashLog(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerCashListResp{List: make([]*types.MoveCashDetail, 0)}
	_ = copier.Copy(resp, rs)
	return
}
