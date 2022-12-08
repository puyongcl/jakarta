package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpireVoiceChatOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExpireVoiceChatOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExpireVoiceChatOrderLogic {
	return &GetExpireVoiceChatOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取需要更新过期的语音订单
func (l *GetExpireVoiceChatOrderLogic) GetExpireVoiceChatOrder(in *pb.GetExpireVoiceChatOrderReq) (*pb.GetExpireVoiceChatOrderResp, error) {
	endExpiryTime, err := time.ParseInLocation(db.DateTimeFormat, in.EndExpiryTime, time.Local)
	if err != nil {
		return nil, err
	}
	expiryTimeStart, err := time.ParseInLocation(db.DateTimeFormat, in.StartExpiryTime, time.Local)
	if err != nil {
		return nil, err
	}
	rsp, err := l.svcCtx.ChatOrderModel.FindExpireOrder(l.ctx, &expiryTimeStart, &endExpiryTime, orderkey.ListenerOrderTypeVoiceChat, orderkey.CanExpiryOrderSate, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		return &pb.GetExpireVoiceChatOrderResp{List: make([]*pb.ExpireVoiceChatOrder, 0)}, nil
	}
	resp := &pb.GetExpireVoiceChatOrderResp{List: make([]*pb.ExpireVoiceChatOrder, 0)}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.ExpireVoiceChatOrder
		_ = copier.Copy(&val, rsp[idx])
		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
