package stat

import (
	"context"
	pbStat "jakarta/app/statistic/rpc/pb"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAdultQuizECRLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAdultQuizECRLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAdultQuizECRLogic {
	return &SaveAdultQuizECRLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAdultQuizECRLogic) SaveAdultQuizECR(req *types.SaveAdultQuizECRReq) (resp *types.SaveAdultQuizECRResp, err error) {
	if req.Uid <= 0 || req.Result <= 0 || len(req.Answer) != 36 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	var in pbStat.SaveAdultQuizEcrReq
	in.Result = req.Result
	in.Answer = req.Answer
	in.Uid = req.Uid
	_, err = l.svcCtx.StatRpc.SaveAdultQuizECR(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.SaveAdultQuizECRResp{}
	return
}
