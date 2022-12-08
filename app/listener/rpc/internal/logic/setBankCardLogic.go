package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	listenerPgModel "jakarta/app/pgModel/listenerPgModel"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetBankCardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBankCardLogic {
	return &SetBankCardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  绑定银行卡
func (l *SetBankCardLogic) SetBankCard(in *pb.SetBankCardReq) (*pb.SetBankCardResp, error) {
	var data *listenerPgModel.ListenerBankCard
	data, err := l.svcCtx.ListenerBankCardModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	if data == nil {
		data = new(listenerPgModel.ListenerBankCard)
		_ = copier.Copy(data, in)

		var pf *listenerPgModel.ListenerProfile
		pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
		if err != nil {
			return nil, err
		}
		data.ListenerName = pf.ListenerName
		data.PhoneNumber = pf.PhoneNumber
		data.IdNo = pf.IdNo

		_, err = l.svcCtx.ListenerBankCardModel.Insert(l.ctx, data)
		if err != nil {
			return nil, err
		}
		return &pb.SetBankCardResp{}, nil
	}
	return l.updateBankCard(in)
}

func (l *SetBankCardLogic) updateBankCard(in *pb.SetBankCardReq) (*pb.SetBankCardResp, error) {
	newData := new(listenerPgModel.ListenerBankCard)
	if in.BankCardNo == "" {
		return &pb.SetBankCardResp{}, nil
	}
	newData.BankCardNo = in.BankCardNo
	newData.ListenerUid = in.ListenerUid

	var pf *listenerPgModel.ListenerProfile
	var err error
	pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil {
		return nil, err
	}
	newData.ListenerName = pf.ListenerName
	newData.PhoneNumber = pf.PhoneNumber
	newData.IdNo = pf.IdNo

	err = l.svcCtx.ListenerBankCardModel.UpdateListenerBankCard(l.ctx, newData)
	if err != nil {
		return nil, err
	}
	return &pb.SetBankCardResp{}, nil
}
