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

type GetListenerWordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListenerWordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerWordsLogic {
	return &GetListenerWordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListenerWordsLogic) GetListenerWords(req *types.GetListenerWordsReq) (resp *types.GetListenerWordsResp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if req.ListenerUid != uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.ListenerUid))
	}
	var in pb.GetListenerWordsReq
	in.ListenerUid = req.ListenerUid
	rs, err := l.svcCtx.ListenerRpc.GetListenerWords(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetListenerWordsResp{WordsSort: make([]int64, 0)}
	_ = copier.Copy(resp, rs)
	return
}
