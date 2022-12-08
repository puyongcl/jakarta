package logic

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/cservice"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/kqueue"
	"jakarta/common/tool"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendListenerListByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecommendListenerListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendListenerListByUserLogic {
	return &GetRecommendListenerListByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  用户获取XXX推荐列表
func (l *GetRecommendListenerListByUserLogic) GetRecommendListenerListByUser(in *pb.GetRecommendListenerByUserReq) (*pb.GetRecommendListenerByUserResp, error) {
	// 获取黑名单
	blk, err := l.svcCtx.UserRedis.GetBlacklist(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	rsp := make([]*pb.UserSeeRecommendListenerProfile, 0)
	resp := &pb.GetRecommendListenerByUserResp{Listener: rsp}

	// 查询休息中 返回空
	if tool.IsInt64ArrayExist(listenerkey.ListenerWorkStateRestingManual, in.WorkState) || tool.IsInt64ArrayExist(listenerkey.ListenerWorkStateRestingAuto, in.WorkState) {
		return resp, nil
	}

	// 统计推荐曝光
	msg := kqueue.UpdateListenerUserStatMessage{
		Time:        time.Now().Format(db.DateTimeFormat),
		Event:       listenerkey.ListenerUserEventRecommend,
		Uid:         in.Uid,
		ListenerUid: make([]int64, 0),
	}
	// 默认可接单
	in.WorkState = []int64{listenerkey.ListenerWorkStateWorking}

	rs1, err := l.svcCtx.ListenerProfileModel.FindRecommendListenerList(l.ctx, in.PageNo, in.PageSize, in.Specialties, in.ChatType, in.Gender, in.Age, in.WorkState, []int64{}, in.SortOrder, blk)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetRecommendListenerListByUserLogic ListenerProfileModel.FindRecommendListenerList err:%+v", err)
		return resp, err
	}
	//
	for idx := 0; idx < len(rs1); idx++ {
		if tool.IsInt64ArrayExist(rs1[idx].ListenerUid, msg.ListenerUid) {
			continue
		}

		if tool.IsStringArrayExist(in.AuthKey, listenerkey.TestUserAuthKey) && rs1[idx].ListenerUid == cservice.DefaultListenerUid {
			rs1[idx].VoiceChatPrice = 100
			rs1[idx].TextChatPrice = 100
		}

		val := addRecListenerData(rs1[idx])
		msg.ListenerUid = append(msg.ListenerUid, rs1[idx].ListenerUid)
		rsp = append(rsp, val)
	}

	if len(msg.ListenerUid) <= 0 && in.PageNo == 1 {
		// 没有查询到符合要求的用户 去掉所有筛选条件
		in.WorkState = []int64{listenerkey.ListenerWorkStateRestingAuto, listenerkey.ListenerWorkStateRestingManual}

		rs1, err = l.svcCtx.ListenerProfileModel.FindRecommendListenerList(l.ctx, in.PageNo, in.PageSize, 0, 0, 0, 0, in.WorkState, []int64{}, in.SortOrder, blk)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GetRecommendListenerListByUserLogic ListenerProfileModel.FindRecommendListenerList err:%+v", err)
			return resp, err
		}
		//
		for idx := 0; idx < len(rs1); idx++ {
			if tool.IsInt64ArrayExist(rs1[idx].ListenerUid, msg.ListenerUid) {
				continue
			}

			if tool.IsStringArrayExist(in.AuthKey, listenerkey.TestUserAuthKey) && rs1[idx].ListenerUid == cservice.DefaultListenerUid {
				rs1[idx].VoiceChatPrice = 100
				rs1[idx].TextChatPrice = 100
			}

			val := addRecListenerData(rs1[idx])
			msg.ListenerUid = append(msg.ListenerUid, rs1[idx].ListenerUid)
			rsp = append(rsp, val)
		}
	}

	l.kqPush(&msg)
	return &pb.GetRecommendListenerByUserResp{Listener: rsp}, nil
}

func addRecListenerData(va *listenerPgModel.ListenerProfile) *pb.UserSeeRecommendListenerProfile {
	var val pb.UserSeeRecommendListenerProfile
	_ = copier.Copy(&val, va)
	val.Age = tool.GetAge2(va.Birthday)
	return &val
}

func (l *GetRecommendListenerListByUserLogic) kqPush(msg *kqueue.UpdateListenerUserStatMessage) {
	if len(msg.ListenerUid) <= 0 {
		return
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetRecommendListenerListByUserLogic kqPush json marshal err:%+v", err)
		return
	}

	err = l.svcCtx.KqUpdateListenerUserStatClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GetRecommendListenerListByUserLogic kqPush Push err:%+v", err)
		return
	}
}
