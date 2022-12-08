package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/tool"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerBasicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerBasicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerBasicInfoLogic {
	return &GetListenerBasicInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX基本资料
func (l *GetListenerBasicInfoLogic) GetListenerBasicInfo(in *pb.GetListenerBasicInfoReq) (*pb.GetListenerBasicInfoResp, error) {
	rs, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	resp := pb.GetListenerBasicInfoResp{}
	if rs != nil {
		_ = copier.Copy(&resp, rs)
	}
	resp.Age = tool.GetAge2(rs.Birthday)
	return &resp, nil
}
