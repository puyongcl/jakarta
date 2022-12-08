package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBankCardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBankCardLogic {
	return &GetBankCardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取银行卡
func (l *GetBankCardLogic) GetBankCard(in *pb.GetBankCardReq) (*pb.GetBankCardResp, error) {
	data, err := l.svcCtx.ListenerBankCardModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	resp := &pb.GetBankCardResp{}
	if data != nil {
		_ = copier.Copy(resp, data)
	}

	// 如果为空，则设置默认值
	if resp.ListenerName == "" || resp.IdNo == "" || resp.PhoneNumber == "" {
		var pf *listenerPgModel.ListenerProfile
		pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil && err != listenerPgModel.ErrNotFound {
			return nil, err
		}
		if pf != nil {
			resp.ListenerName = pf.ListenerName
			resp.PhoneNumber = pf.PhoneNumber
			resp.IdNo = pf.IdNo
		}
	}

	return resp, nil
}
