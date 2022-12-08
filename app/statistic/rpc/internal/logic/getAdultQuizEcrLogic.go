package logic

import (
	"context"
	"jakarta/app/pgModel/statPgModel"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdultQuizEcrLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdultQuizEcrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdultQuizEcrLogic {
	return &GetAdultQuizEcrLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取最新成人依恋量表测试结果
func (l *GetAdultQuizEcrLogic) GetAdultQuizEcr(in *pb.GetAdultQuizEcrReq) (*pb.GetAdultQuizEcrResp, error) {
	rs, err := l.svcCtx.AdultQuizEcrModel.FindOneByUid(l.ctx, in.Uid)
	if err != nil && err != statPgModel.ErrNotFound {
		return nil, err
	}

	resp := pb.GetAdultQuizEcrResp{}
	if rs != nil {
		resp.Result = rs.Result
	}

	return &resp, nil
}
