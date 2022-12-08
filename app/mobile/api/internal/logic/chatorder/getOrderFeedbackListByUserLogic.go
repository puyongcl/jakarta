package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderFeedbackListByUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderFeedbackListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderFeedbackListByUserLogic {
	return &GetOrderFeedbackListByUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderFeedbackListByUserLogic) GetOrderFeedbackListByUser(req *types.GetChatOrderFeedbackListByUserReq) (resp *types.GetChatOrderFeedbackListByUserResp, err error) {
	if req.Uid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	var in pb.GetChatOrderFeedbackListByUserReq
	in.Uid = req.Uid
	in.PageNo = req.PageNo
	in.PageSize = req.PageSize
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderFeedbackListByUser(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetChatOrderFeedbackListByUserResp{List: make([]*types.UserSeeChatOrderFeedback, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
