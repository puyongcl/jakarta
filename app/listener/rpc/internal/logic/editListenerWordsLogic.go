package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditListenerWordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditListenerWordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditListenerWordsLogic {
	return &EditListenerWordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  编辑XXX常用语 删除1条常用语 该条后方的向前挪位置 排序数组 去掉末尾序号 排序数组内的数字代表了实际的条数和顺序
func (l *EditListenerWordsLogic) EditListenerWords(in *pb.EditListenerWordsReq) (*pb.EditListenerWordsResp, error) {
	newData := new(listenerPgModel.ListenerWords)
	_ = copier.Copy(newData, in)
	err := l.svcCtx.ListenerWordsModel.UpdatePart(l.ctx, newData)
	if err != nil {
		return nil, err
	}
	return &pb.EditListenerWordsResp{}, nil
}
