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

type GetListenerOwnProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerOwnProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerOwnProfileLogic {
	return &GetListenerOwnProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerOwnProfileLogic) GetListenerOwnProfile(req *types.GetListenerOwnInfoReq) (resp *types.GetListenerOwnInfoResp, err error) {
	return l.GetListenerOwnProfile2(req, true)
}

func (l *GetListenerOwnProfileLogic) GetListenerOwnProfile2(req *types.GetListenerOwnInfoReq, valid bool) (resp *types.GetListenerOwnInfoResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if valid && req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	rs, err := l.svcCtx.ListenerRpc.GetListenerProfileByOwn(l.ctx, &pb.GetListenerProfileByOwnReq{ListenerUid: req.ListenerUid})
	if err != nil {
		return nil, err
	}
	var val types.ListenerSeeOwnProfile
	_ = copier.Copy(&val, rs.Profile)
	resp = &types.GetListenerOwnInfoResp{Info: &val}
	return
}
