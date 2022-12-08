package logic

import (
	"context"
	"jakarta/common/key/db"
	"time"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserChannelListLogic {
	return &GetUserChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取用户渠道列表
func (l *GetUserChannelListLogic) GetUserChannelList(in *pb.GetUserChannelListReq) (*pb.GetUserChannelListResp, error) {
	var start, end time.Time
	var err error
	if in.CreateTimeStart != "" {
		start, err = time.ParseInLocation(db.DateTimeFormat, in.CreateTimeStart, time.Local)
		if err != nil {
			return nil, err
		}
	}
	if in.CreateTimeEnd != "" {
		end, err = time.ParseInLocation(db.DateTimeFormat, in.CreateTimeEnd, time.Local)
		if err != nil {
			return nil, err
		}
	}
	resp := pb.GetUserChannelListResp{Channel: make([]string, 0)}

	resp.Channel, err = l.svcCtx.UserAuthModel.FindChannelList(l.ctx, &start, &end)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
