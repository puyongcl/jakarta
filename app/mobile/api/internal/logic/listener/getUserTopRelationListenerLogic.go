package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"sync"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTopRelationListenerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserTopRelationListenerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTopRelationListenerLogic {
	return &GetUserTopRelationListenerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserTopRelationListenerLogic) GetUserTopRelationListener(req *types.GetUserTopRelationListenerReq) (resp *types.GetUserTopRelationListenerResp, err error) {
	in1 := pbChat.GetTopUserAndListenerRelationReq{
		Uid:      req.Uid,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}
	rsp1, err := l.svcCtx.ChatRpc.GetTopUserAndListenerRelation(l.ctx, &in1)
	if err != nil {
		return nil, err
	}

	resp = &types.GetUserTopRelationListenerResp{}

	if len(rsp1.List) <= 0 {
		return
	}

	mL := &sync.Map{}
	fns := make([]func() error, 0)
	// 查询
	for idx := 0; idx < len(rsp1.List); idx++ {
		listenerUid := rsp1.List[idx].ListenerUid
		fns = append(fns, func() error {
			var lin pbListener.GetListenerBasicInfoReq
			lin.ListenerUid = listenerUid
			rsp2, err2 := l.svcCtx.ListenerRpc.GetListenerBasicInfo(l.ctx, &lin)
			if err2 != nil {
				return err2
			}
			mL.Store(lin.ListenerUid, rsp2)
			return nil
		})
	}

	err = mr.Finish(fns...)
	if err != nil {
		return nil, err
	}

	//
	resp = &types.GetUserTopRelationListenerResp{List: make([]*types.UserSeeListenerShortProfile, 0)}
	for idx := 0; idx < len(rsp1.List); idx++ {
		var val types.UserSeeListenerShortProfile
		vv, ok := mL.Load(rsp1.List[idx].ListenerUid)
		if ok {
			lbc := vv.(*pbListener.GetListenerBasicInfoResp)
			_ = copier.Copy(&val, lbc)
		}

		resp.List = append(resp.List, &val)
	}
	return
}
