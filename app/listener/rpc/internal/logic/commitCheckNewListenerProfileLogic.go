package logic

import (
	"context"
	"fmt"
	"jakarta/common/key/listenerkey"
	"jakarta/common/xerr"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommitCheckNewListenerProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommitCheckNewListenerProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommitCheckNewListenerProfileLogic {
	return &CommitCheckNewListenerProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  新申请XXX提交审核
func (l *CommitCheckNewListenerProfileLogic) CommitCheckNewListenerProfile(in *pb.CommitCheckNewListenerProfileReq) (*pb.CommitCheckNewListenerProfileResp, error) {
	data, err := l.svcCtx.ListenerProfileDraftModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	if data.CheckStatus != listenerkey.CheckStatusFirstApplyEdit {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("当前状态无法提交审核 %d", data.CheckStatus))
	}
	err = l.svcCtx.ListenerProfileDraftModel.UpdateCheckStatus(l.ctx, in.Uid, in.CheckStatus)
	if err != nil {
		return nil, err
	}
	return &pb.CommitCheckNewListenerProfileResp{}, nil
}
