package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerWordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerWordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerWordsLogic {
	return &GetListenerWordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX常用语
func (l *GetListenerWordsLogic) GetListenerWords(in *pb.GetListenerWordsReq) (*pb.GetListenerWordsResp, error) {
	rs, err := l.svcCtx.ListenerWordsModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetListenerWordsResp{}
	_ = copier.Copy(resp, rs)
	return resp, nil
}
