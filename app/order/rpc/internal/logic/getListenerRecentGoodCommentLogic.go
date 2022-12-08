package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"strings"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerRecentGoodCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerRecentGoodCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerRecentGoodCommentLogic {
	return &GetListenerRecentGoodCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取指定XXX的好评
func (l *GetListenerRecentGoodCommentLogic) GetListenerRecentGoodComment(in *pb.GetListenerRecentGoodCommentReq) (*pb.GetListenerRecentGoodCommentResp, error) {
	resp := &pb.GetListenerRecentGoodCommentResp{List: make([]*pb.RecentGoodComment, 0)}

	for idx := 0; idx < len(in.Listener); idx++ {
		rsp, err := l.svcCtx.ChatOrderModel.FindLastCommentOrder(l.ctx, in.Listener[idx].ListenerUid, listenerkey.Rating5Star)
		if err != nil && err != orderPgModel.ErrNotFound {
			return nil, err
		}
		if rsp == nil {
			err = nil
			continue
		}
		var txt string
		if idx%2 == 1 {
			txt = fmt.Sprintf(orderkey.TipShowGoodComment, rsp.NickName, strings.Replace(in.Listener[idx].ListenerNickName, "#", "", -1))
		} else {
			txt = fmt.Sprintf(orderkey.TipShowServiceOut, strings.Replace(in.Listener[idx].ListenerNickName, "#", "", -1), rsp.NickName)
		}
		var val pb.RecentGoodComment
		val.Text = txt
		val.Uid = rsp.Uid
		val.ListenerUid = rsp.ListenerUid
		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
