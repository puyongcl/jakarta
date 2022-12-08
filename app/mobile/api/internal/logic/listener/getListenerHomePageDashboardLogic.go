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

type GetListenerHomePageDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerHomePageDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerHomePageDashboardLogic {
	return &GetListenerHomePageDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerHomePageDashboardLogic) GetListenerHomePageDashboard(req *types.GetListenerHomePageDashboardReq) (resp *types.GetListenerHomePageDashboardResp, err error) {
	return l.GetListenerHomePageDashboard2(req, true)
}

func (l *GetListenerHomePageDashboardLogic) GetListenerHomePageDashboard2(req *types.GetListenerHomePageDashboardReq, valid bool) (resp *types.GetListenerHomePageDashboardResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	if valid && req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	rs, err := l.svcCtx.ListenerRpc.GetListenerHomePageDashboard(l.ctx, &pbListener.GetListenerHomePageDashboardReq{ListenerUid: req.ListenerUid})
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerHomePageDashboardResp{}
	if rs != nil {
		_ = copier.Copy(resp, rs)
	}
	return
}
