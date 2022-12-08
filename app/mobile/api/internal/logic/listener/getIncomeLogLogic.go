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

type GetIncomeLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncomeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncomeLogLogic {
	return &GetIncomeLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncomeLogLogic) GetIncomeLog(req *types.GetListenerIncomeListReq) (resp *types.GetListenerIncomeListResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	var in pbListener.GetListenerIncomeLogReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.ListenerRpc.GetListenerIncomeLog(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerIncomeListResp{List: make([]*types.ListenerIncomeDetail, 0)}
	_ = copier.Copy(resp, rs)
	return
}
