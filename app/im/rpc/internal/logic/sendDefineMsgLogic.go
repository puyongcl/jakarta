package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"
	"jakarta/common/third_party/tim"
)

type SendDefineMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendDefineMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendDefineMsgLogic {
	return &SendDefineMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发送自定义消息
func (l *SendDefineMsgLogic) SendDefineMsg(in *pb.SendDefineMsgReq) (*pb.SendDefineMsgResp, error) {
	mc := tim.DefineMsgContent{
		DefineMsgType: in.MsgType,
		Title:         in.Title,
		Text:          in.Text,
		Val1:          in.Val1,
		Val2:          in.Val2,
		Val3:          in.Val3,
		Val4:          in.Val4,
		Val5:          in.Val5,
		Val6:          in.Val6,
	}
	rs, err := l.svcCtx.TimClient.SendDefineMsg(in.FromUid, in.ToUid, &mc, in.Sync)
	if err != nil {
		return nil, err
	}
	return &pb.SendDefineMsgResp{ErrCode: rs.ErrorCode, ErrMsg: rs.ErrorInfo}, nil
}
