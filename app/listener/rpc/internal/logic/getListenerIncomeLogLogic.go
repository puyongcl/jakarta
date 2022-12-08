package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerIncomeLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerIncomeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerIncomeLogLogic {
	return &GetListenerIncomeLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取收益记录
func (l *GetListenerIncomeLogLogic) GetListenerIncomeLog(in *pb.GetListenerIncomeLogReq) (*pb.GetListenerIncomeLogResp, error) {
	rs, err := l.svcCtx.ListenerWalletFlowModel.Find(l.ctx, in.ListenerUid, listenerkey.IncomeSettleType, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetListenerIncomeLogResp{List: make([]*pb.ListenerIncomeLog, 0)}
	if len(rs) <= 0 {
		return resp, nil
	}
	for idx := 0; idx < len(rs); idx++ {
		var val pb.ListenerIncomeLog
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].OutTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
