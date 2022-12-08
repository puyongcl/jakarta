package bbs

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"jakarta/app/bbs/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/xerr"
	"sync"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoryReplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStoryReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoryReplyListLogic {
	return &GetStoryReplyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type BasicInfo struct {
	Uid      int64
	NickName string
	Avatar   string
	Intro    string
}

func (l *GetStoryReplyListLogic) GetStoryReplyList(req *types.GetStoryReplyListReq) (resp *types.GetStoryReplyListResp, err error) {
	if req.StoryId == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "Uid参数为空")
	}
	var in pb.GetStoryReplyListByUserReq
	_ = copier.Copy(&in, req)
	in.Uid = uid

	rsp, err := l.svcCtx.BbsRpc.GetStoryReplyListByUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}

	if len(rsp.List) <= 0 {
		return &types.GetStoryReplyListResp{}, nil
	}

	mL := &sync.Map{}
	fns := make([]func() error, 0)
	// 查询
	for idx := 0; idx < len(rsp.List); idx++ {
		listenerUid := rsp.List[idx].ListenerUid
		fns = append(fns, func() error {
			var lin pbListener.GetListenerBasicInfoReq
			lin.ListenerUid = listenerUid
			rsp1, err2 := l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &lin)
			if err2 != nil {
				return err2
			}
			mL.Store(lin.ListenerUid, &BasicInfo{
				Uid:      rsp1.ListenerUid,
				NickName: rsp1.NickName,
				Avatar:   rsp1.Avatar,
				Intro:    rsp1.Introduction,
			})
			return nil
		})
	}

	err = mr.Finish(fns...)
	if err != nil {
		return nil, err
	}

	//
	resp = &types.GetStoryReplyListResp{List: make([]*types.StoryReply, 0)}
	for idx := 0; idx < len(rsp.List); idx++ {
		var val types.StoryReply
		_ = copier.Copy(&val, rsp.List[idx])
		vv, ok := mL.Load(val.ListenerUid)
		if ok {
			lbc := vv.(*BasicInfo)
			val.Intro = lbc.Intro
			val.ListenerAvatar = lbc.Avatar
			val.ListenerNickname = lbc.NickName
		}

		resp.List = append(resp.List, &val)
	}

	return
}
