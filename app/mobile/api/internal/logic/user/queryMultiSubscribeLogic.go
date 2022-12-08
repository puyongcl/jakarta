package user

import (
	"context"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryMultiSubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryMultiSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMultiSubscribeLogic {
	return &QueryMultiSubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryMultiSubscribeLogic) QueryMultiSubscribe(req *types.QueryMultiSubscribeNotifyMsgReq) (resp *types.QueryMultiSubscribeNotifyMsgResp, err error) {
	if len(req.Sub) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}

	resp = &types.QueryMultiSubscribeNotifyMsgResp{Rsp: make([]*types.QuerySubscribeNotifyMsgResp, 0)}
	var b bool
	for k, _ := range req.Sub {
		// 判断用户是否有订阅
		b, err = l.svcCtx.IMRedis.IsUserHaveOneTimeSubscribeMsg(l.ctx, req.Sub[k].MsgType, req.Sub[k].Uid)
		if err != nil {
			return
		}
		rsp := &types.QuerySubscribeNotifyMsgResp{
			MsgType: req.Sub[k].MsgType,
			SendCnt: 0,
		}
		if b {
			rsp.SendCnt = 1
		}

		resp.Rsp = append(resp.Rsp, rsp)
	}

	return
}
