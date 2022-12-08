package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatOrderFeedbackListByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatOrderFeedbackListByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatOrderFeedbackListByUserLogic {
	return &GetChatOrderFeedbackListByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户获取反馈列表
func (l *GetChatOrderFeedbackListByUserLogic) GetChatOrderFeedbackListByUser(in *pb.GetChatOrderFeedbackListByUserReq) (*pb.GetChatOrderFeedbackListByUserResp, error) {
	rsp, err := l.svcCtx.ChatOrderModel.FindChatOrderFeedback(l.ctx, in.Uid, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		return &pb.GetChatOrderFeedbackListByUserResp{}, nil
	}
	resp := pb.GetChatOrderFeedbackListByUserResp{List: make([]*pb.UserSeeChatOrderFeedback, 0)}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.UserSeeChatOrderFeedback
		_ = copier.Copy(&val, rsp[idx])
		val.CreateTime = rsp[idx].CreateTime.Format(db.DateTimeFormat)
		if rsp[idx].FeedbackTime.Valid {
			val.FeedbackTime = rsp[idx].FeedbackTime.Time.Format(db.DateFormat)
		}

		resp.List = append(resp.List, &val)
	}
	return &resp, nil
}
