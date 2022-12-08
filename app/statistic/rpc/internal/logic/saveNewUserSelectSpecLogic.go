package logic

import (
	"context"
	"jakarta/app/pgModel/statPgModel"
	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNewUserSelectSpecLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNewUserSelectSpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNewUserSelectSpecLogic {
	return &SaveNewUserSelectSpecLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  新用户选择的XX标签
func (l *SaveNewUserSelectSpecLogic) SaveNewUserSelectSpec(in *pb.SaveNewUserSelectSpecReq) (*pb.SaveNewUserSelectSpecResp, error) {
	data := &statPgModel.UserSelectSpecTag{
		Uid:         in.Uid,
		Specialties: in.Spec,
		Channel:     in.Channel,
	}

	_, err := l.svcCtx.UserSelectSpecTagModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &pb.SaveNewUserSelectSpecResp{}, nil
}
