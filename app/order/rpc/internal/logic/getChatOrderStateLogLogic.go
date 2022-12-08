package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatOrderStateLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatOrderStateLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatOrderStateLogLogic {
	return &GetChatOrderStateLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取用户订单状态变化记录
func (l *GetChatOrderStateLogLogic) GetChatOrderStateLog(in *pb.GetChatOrderStateLogReq) (*pb.GetChatOrderStateLogResp, error) {
	rsp, err := l.svcCtx.ChatOrderStatusLogModel.Find(l.ctx, in.OrderId, in.State, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := pb.GetChatOrderStateLogResp{List: make([]*pb.ChatOrderStateLog, 0)}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.ChatOrderStateLog
		_ = copier.Copy(&val, rsp[idx])
		val.CreateTime = rsp[idx].CreateTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}
	return &resp, nil
}
