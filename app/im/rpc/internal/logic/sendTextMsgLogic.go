package logic

import (
	"context"
	"github.com/jinzhu/copier"

	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTextMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendTextMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTextMsgLogic {
	return &SendTextMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发送普通消息
func (l *SendTextMsgLogic) SendTextMsg(in *pb.SendTextMsgReq) (*pb.SendTextMsgResp, error) {
	rsp, err := l.svcCtx.TimClient.SendTextMsg(in.FromUid, in.ToUid, in.Text, in.Sync)
	if err != nil {
		return nil, err
	}
	resp := &pb.SendTextMsgResp{}
	_ = copier.Copy(resp, rsp)
	return resp, nil
}
