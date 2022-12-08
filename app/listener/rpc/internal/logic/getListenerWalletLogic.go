package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerWalletLogic {
	return &GetListenerWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX钱包详情
func (l *GetListenerWalletLogic) GetListenerWallet(in *pb.GetListenerWalletReq) (*pb.GetListenerWalletResp, error) {
	data, err := l.svcCtx.ListenerWalletModel.FindOne2(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	resp := &pb.GetListenerWalletResp{}
	if data != nil {
		_ = copier.Copy(resp, data)
	}

	return resp, nil
}
