package user

import (
	"context"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/notify"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerySubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuerySubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerySubscribeLogic {
	return &QuerySubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerySubscribeLogic) QuerySubscribe(req *types.QuerySubscribeNotifyMsgReq) (resp *types.QuerySubscribeNotifyMsgResp, err error) {
	switch req.MsgType {
	case notify.DefineNotifyMsgTypeMiniProgramMsg1:
		return l.oneTime(req)
	}
	return
}

func (l *QuerySubscribeLogic) oneTime(req *types.QuerySubscribeNotifyMsgReq) (resp *types.QuerySubscribeNotifyMsgResp, err error) {
	// 判断用户是否有订阅
	b, err := l.svcCtx.IMRedis.IsUserHaveOneTimeSubscribeMsg(l.ctx, req.MsgType, req.Uid)
	if err != nil {
		return
	}
	resp = &types.QuerySubscribeNotifyMsgResp{
		MsgType: req.MsgType,
		SendCnt: 0,
	}
	if !b {
		return
	}
	resp.SendCnt = 1
	return
}
