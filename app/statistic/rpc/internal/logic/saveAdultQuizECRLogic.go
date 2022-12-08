package logic

import (
	"context"
	"jakarta/app/pgModel/statPgModel"
	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"
	"jakarta/common/uniqueid"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAdultQuizECRLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveAdultQuizECRLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAdultQuizECRLogic {
	return &SaveAdultQuizECRLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存成人依恋量表测试结果
func (l *SaveAdultQuizECRLogic) SaveAdultQuizECR(in *pb.SaveAdultQuizEcrReq) (*pb.SaveAdultQuizEcrResp, error) {
	data := statPgModel.AdultQuizEcr{
		Id:     uniqueid.GenDataId(),
		Uid:    in.Uid,
		Result: in.Result,
		Answer: in.Answer,
	}

	_, err := l.svcCtx.AdultQuizEcrModel.Insert(l.ctx, &data)
	if err != nil {
		return nil, err
	}

	return &pb.SaveAdultQuizEcrResp{}, nil
}
