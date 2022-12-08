package logic

import (
	"context"
	"jakarta/common/key/db"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindListenerListRangeByUpdateTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindListenerListRangeByUpdateTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindListenerListRangeByUpdateTimeLogic {
	return &FindListenerListRangeByUpdateTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  查询几天内更新过的XXX列表
func (l *FindListenerListRangeByUpdateTimeLogic) FindListenerListRangeByUpdateTime(in *pb.FindListenerListRangeByUpdateTimeReq) (*pb.FindListenerListRangeByUpdateTimeResp, error) {
	start, err := time.ParseInLocation(db.DateTimeFormat, in.Start, time.Local)
	if err != nil {
		return nil, err
	}
	end, err := time.ParseInLocation(db.DateTimeFormat, in.End, time.Local)
	if err != nil {
		return nil, err
	}
	resp := &pb.FindListenerListRangeByUpdateTimeResp{Listener: make([]*pb.ListenerShortProfile, 0)}
	rsp, err := l.svcCtx.ListenerProfileModel.FindListenerUidRangeUpdateTime(l.ctx, in.PageNo, in.PageSize, in.WorkState, &start, &end)
	if err != nil {
		return nil, err
	}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.ListenerShortProfile
		val.ListenerUid = rsp[idx].ListenerUid
		val.ListenerNickName = rsp[idx].NickName
		val.ListenerAvatar = rsp[idx].Avatar

		resp.Listener = append(resp.Listener, &val)
	}

	return resp, nil
}
