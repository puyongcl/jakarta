package listener

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRemarkUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerRemarkUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRemarkUserLogic {
	return &GetListenerRemarkUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerRemarkUserLogic) GetListenerRemarkUser(req *types.GetListenerRemarkUserReq) (resp *types.GetListenerRemarkUserResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	rs, err := l.svcCtx.ListenerRpc.GetListenerRemarkUser(l.ctx, &pb.GetListenerRemarkUserReq{
		Uid:         req.Uid,
		ListenerUid: req.ListenerUid,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerRemarkUserResp{}
	_ = copier.Copy(&resp, rs)
	return
}
