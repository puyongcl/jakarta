package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/admin/api/internal/logic/adminlog"
	pbChat "jakarta/app/chat/rpc/pb"
	pbBusiness "jakarta/app/listener/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/userkey"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckListenerProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckListenerProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckListenerProfileLogic {
	return &CheckListenerProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckListenerProfileLogic) CheckListenerProfile(req *types.CheckListenerProfileReq) (resp *types.CheckListenerProfileResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "CheckListenerProfile", req.AdminUid, err, req, resp)
	}()
	var in pbBusiness.CheckListenerProfileReq
	// 首次申请通过取openid
	if req.CheckStatus == listenerkey.CheckStatusFirstApplyChecking && len(req.CheckFailField) == 0 {
		var in3 pbUser.GetUserAuthByUserIdReq
		in3.Uid = req.ListenerUid
		var rs *pbUser.GetUserAuthyUserIdResp
		rs, err = l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &in3)
		if err != nil {
			return nil, err
		}
		in.OpenId = rs.UserAuth.AuthKey
		in.Channel = rs.UserAuth.Channel
	}
	_ = copier.Copy(&in, req)
	rsp, err := l.svcCtx.ListenerRpc.AdminCheckListenerProfile(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.CheckListenerProfileResp{}
	//_ = copier.Copy(resp, rsp)

	// 更新用户类型
	if rsp.CheckResult == listenerkey.CheckStatusFirstApplyPass {
		in2 := pbUser.UpdateUserTypeReq{
			Uid:      req.ListenerUid,
			UserType: userkey.UserTypeListener,
		}
		_, err = l.svcCtx.UsercenterRpc.UpdateUserType(l.ctx, &in2)
		if err != nil {
			return nil, err
		}

		// 初始化XXX通话状态
		in3 := pbChat.CreateListenerVoiceChatStateReq{ListenerUid: req.ListenerUid}
		_, err = l.svcCtx.ChatRpc.CreateListenerVoiceChatState(l.ctx, &in3)
		if err != nil {
			return nil, err
		}
	}

	return
}
