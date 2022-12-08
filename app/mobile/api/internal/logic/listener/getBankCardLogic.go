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

type GetBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBankCardLogic {
	return &GetBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBankCardLogic) GetBankCard(req *types.GetBankCardReq) (resp *types.GetBankCardResp, err error) {
	if req.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "empty uid")
	}
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	var in pbListener.GetBankCardReq
	in.ListenerUid = req.ListenerUid
	rs, err := l.svcCtx.ListenerRpc.GetBankCard(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetBankCardResp{}
	_ = copier.Copy(resp, rs)
	return
}
