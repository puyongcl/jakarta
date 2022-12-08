package logic

import (
	"context"
	"fmt"
	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"
	"jakarta/common/key/orderkey"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentGoodCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecentGoodCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecentGoodCommentLogic {
	return &GetRecentGoodCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取最近的好评
func (l *GetRecentGoodCommentLogic) GetRecentGoodComment(in *pb.GetRecentGoodCommentReq) (*pb.GetRecentGoodCommentResp, error) {
	rsp, err := l.svcCtx.ChatOrderModel.FindGoodComment(l.ctx, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetRecentGoodCommentResp{List: make([]*pb.RecentGoodComment, 0)}

	for idx := 0; idx < len(rsp); idx++ {
		var txt string
		if idx%2 == 1 {
			txt = fmt.Sprintf(orderkey.TipShowGoodComment, rsp[idx].NickName, strings.Replace(rsp[idx].ListenerNickName, "#", "", -1))
		} else {
			txt = fmt.Sprintf(orderkey.TipShowServiceOut, strings.Replace(rsp[idx].ListenerNickName, "#", "", -1), rsp[idx].NickName)
		}
		var val pb.RecentGoodComment
		val.Text = txt
		val.Uid = rsp[idx].Uid
		val.ListenerUid = rsp[idx].ListenerUid
		resp.List = append(resp.List, &val)
	}
	return resp, nil
}
