package stat

import (
	"context"
	pbStat "jakarta/app/statistic/rpc/pb"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNewUserSelectSpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveNewUserSelectSpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNewUserSelectSpecLogic {
	return &SaveNewUserSelectSpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveNewUserSelectSpecLogic) SaveNewUserSelectSpec(req *types.SaveNewUserSelectSpecReq) (resp *types.SaveNewUserSelectSpecResp, err error) {
	var in pbStat.SaveNewUserSelectSpecReq
	in.Spec = req.Spec
	in.Uid = req.Uid
	in.Channel = req.Channel
	_, err = l.svcCtx.StatRpc.SaveNewUserSelectSpec(l.ctx, &in)
	if err != nil {
		return nil, err
	}

	return
}
