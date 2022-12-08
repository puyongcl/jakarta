package stat

import (
	"context"
	pbStat "jakarta/app/statistic/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdultQuizECRLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdultQuizECRLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdultQuizECRLogic {
	return &GetAdultQuizECRLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdultQuizECRLogic) GetAdultQuizECR(req *types.GetAdultQuizECRReq) (resp *types.GetAdultQuizECRResp, err error) {
	var in pbStat.GetAdultQuizEcrReq
	in.Uid = req.Uid
	rs, err := l.svcCtx.StatRpc.GetAdultQuizEcr(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetAdultQuizECRResp{Result: rs.Result}
	return
}
