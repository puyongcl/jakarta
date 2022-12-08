package logic

import (
	"context"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/listenerkey"
	"jakarta/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendListenerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecommendListenerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendListenerListLogic {
	return &GetRecommendListenerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取推荐的XXX
func (l *GetRecommendListenerListLogic) GetRecommendListenerList(in *pb.GetRecommendListenerReq) (*pb.GetRecommendListenerResp, error) {
	// 获取黑名单
	blk, err := l.svcCtx.UserRedis.GetBlacklist(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	rsp := make([]*pb.ListenerShortProfile, 0)
	resp := &pb.GetRecommendListenerResp{Listener: rsp}
	var listenerUids []int64
	// 可接单
	rs1, err := l.svcCtx.ListenerProfileModel.FindRecommendListenerList(l.ctx, in.PageNo, in.PageSize, in.Specialties, in.ChatType, in.Gender, in.Age, []int64{listenerkey.ListenerWorkStateWorking}, []int64{}, in.SortOrder, blk)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetRecommendListenerList ListenerProfileModel.FindRecommendListenerList err:%+v", err)
		return resp, err
	}
	//
	for idx := 0; idx < len(rs1); idx++ {
		var val pb.ListenerShortProfile
		val.ListenerUid = rs1[idx].ListenerUid
		val.ListenerNickName = rs1[idx].NickName
		val.ListenerAvatar = rs1[idx].Avatar
		listenerUids = append(listenerUids, rs1[idx].ListenerUid)

		rsp = append(rsp, &val)
	}

	// 去掉 可接单条件
	if len(rs1) <= 0 && in.PageNo <= 1 {
		var rs2 []*listenerPgModel.ListenerProfile
		rs2, err = l.svcCtx.ListenerProfileModel.FindRecommendListenerList(l.ctx, in.PageNo, in.PageSize, in.Specialties, in.ChatType, in.Gender, in.Age, []int64{}, []int64{}, in.SortOrder, blk)
		if err != nil {
			return nil, err
		}
		//
		for idx := 0; idx < len(rs2); idx++ {
			if tool.IsInt64ArrayExist(rs2[idx].ListenerUid, listenerUids) {
				continue
			}
			var val pb.ListenerShortProfile
			val.ListenerUid = rs2[idx].ListenerUid
			val.ListenerNickName = rs2[idx].NickName
			val.ListenerAvatar = rs2[idx].Avatar
			listenerUids = append(listenerUids, rs2[idx].ListenerUid)

			rsp = append(rsp, &val)
		}
	}
	return &pb.GetRecommendListenerResp{Listener: rsp}, nil
}
