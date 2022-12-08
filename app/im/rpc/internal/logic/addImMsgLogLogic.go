package logic

import (
	"context"
	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"
	"jakarta/app/pgModel/imPgModel"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddImMsgLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddImMsgLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddImMsgLogLogic {
	return &AddImMsgLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  im msg log
func (l *AddImMsgLogLogic) AddImMsgLog(in *pb.AddImMsgLogReq) (*pb.AddImMsgLogResp, error) {
	mt := time.Now()
	if in.MsgTime > 0 {
		mt = time.Unix(in.MsgTime, 0)
	}

	data := imPgModel.ImMsgLog{
		FromUid:      in.FromUid,
		ToUid:        in.ToUid,
		MsgTime:      mt,
		MsgId:        in.MsgId,
		MsgType:      in.MsgType,
		MsgSeq:       in.MsgSeq,
		FromUserType: in.FromUserType,
	}

	_, err := l.svcCtx.ImMsgLogModel.Insert(l.ctx, &data)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AddImMsgLogLogic err:%+v", err)
	}
	return &pb.AddImMsgLogResp{}, nil
}
