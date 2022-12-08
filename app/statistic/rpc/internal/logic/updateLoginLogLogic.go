package logic

import (
	"context"
	"fmt"
	"jakarta/common/key/db"
	"time"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoginLogLogic {
	return &UpdateLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新每日登陆时间
func (l *UpdateLoginLogLogic) UpdateLoginLog(in *pb.UpdateLoginLogReq) (*pb.UpdateLoginLogResp, error) {
	id := fmt.Sprintf("%d-%s", in.Uid, time.Now().Format(db.DateFormat2))
	_, err := l.svcCtx.UserLoginLogModel.InsertOrUpdate(l.ctx, id, in.Uid, in.UserType, in.Channel)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateLoginLogResp{}, nil
}
