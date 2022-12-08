package logic

import (
	"context"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/cservice"
	"jakarta/common/key/listenerkey"
	"jakarta/common/tool"
	"jakarta/common/xerr"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerPriceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerPriceLogic {
	return &GetListenerPriceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取XXX定价和价格方案
func (l *GetListenerPriceLogic) GetListenerPrice(in *pb.GetListenerPriceReq) (*pb.GetListenerPriceResp, error) {
	rs, err := l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	if rs == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "listener not find")
	}

	if tool.IsStringArrayExist(in.AuthKey, listenerkey.TestUserAuthKey) && rs.ListenerUid == cservice.DefaultListenerUid {
		rs.VoiceChatPrice = 100
		rs.TextChatPrice = 100
	}

	return &pb.GetListenerPriceResp{
		TextChatPrice:  rs.TextChatPrice,
		VoiceChatPrice: rs.VoiceChatPrice,
		Channel:        rs.Channel,
	}, nil
}
