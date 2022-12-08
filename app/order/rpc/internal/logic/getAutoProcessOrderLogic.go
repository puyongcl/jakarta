package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAutoProcessOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAutoProcessOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAutoProcessOrderLogic {
	return &GetAutoProcessOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取需要自动处理的订单
func (l *GetAutoProcessOrderLogic) GetAutoProcessOrder(in *pb.GetAutoProcessOrderReq) (*pb.GetAutoProcessOrderResp, error) {
	if in.BeforeTime == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "time is empty")
	}
	bTime, err := time.ParseInLocation(db.DateTimeFormat, in.BeforeTime, time.Local)
	if err != nil {
		return nil, err
	}
	rs, err := l.svcCtx.ChatOrderModel.FindNeedAutoProcessOrderList(l.ctx, in.State, &bTime, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetAutoProcessOrderResp{List: make([]*pb.AutoProcessOrder, 0)}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.AutoProcessOrder
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].CreateTime.Format(db.DateTimeFormat)
		val.UpdateTime = rs[idx].UpdateTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
