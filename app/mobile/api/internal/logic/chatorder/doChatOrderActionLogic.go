package chatorder

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	pbChat "jakarta/app/chat/rpc/pb"
	"jakarta/app/order/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/key/userkey"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoChatOrderActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDoChatOrderActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoChatOrderActionLogic {
	return &DoChatOrderActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DoChatOrderActionLogic) DoChatOrderAction(req *types.DoChatOrderActionReq) (resp *types.DoChatOrderActionResp, err error) {
	// 加分布式锁
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid != req.OperatorUid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.OperatorUid))
	}
	rkey := fmt.Sprintf(rediskey.RedisLockUser, uid)
	rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	rl.SetExpire(2)
	b, err := rl.AcquireCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	}
	defer func() {
		_, err2 := rl.ReleaseCtx(l.ctx)
		if err2 != nil {
			logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
			return
		}
	}()

	switch req.Action {
	case orderkey.ChatOrderStateApplyRefund5: // 退款
		return l.applyRefund(req)
	case orderkey.ChatOrderStateUserStopService12: // 主动停止服务
		// 查询用户是否在通话中 通话中不能结束订单
		var in pbChat.SyncChatStateReq
		in.Uid = req.OperatorUid
		in.ListenerUid = req.ListenerUid
		in.Action = chatkey.ChatAction12
		var rsp *pbChat.SyncChatStateResp
		rsp, err = l.svcCtx.ChatRpc.SyncChatState(l.ctx, &in)
		if err != nil {
			return nil, err
		}
		if rsp.CurrentVoiceChatState == chatkey.VoiceChatStateStart || rsp.CurrentVoiceChatState == chatkey.VoiceChatStateStop {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.OrderErrorNotAllowStopOrder, "正在通话中，无法结束订单")
		}

		if req.OrderType != 0 && req.OrderId == "" {
			return l.stopAllChatOrder(req)
		}
		return l.defaultAction(req)

	default:
		return l.defaultAction(req)
	}
}

func (l *DoChatOrderActionLogic) stopAllChatOrder(req *types.DoChatOrderActionReq) (resp *types.DoChatOrderActionResp, err error) {
	var in pb.UpdateChatOrderUseReq
	_ = copier.Copy(&in, req)
	in.Uid = req.OperatorUid
	_, err = l.svcCtx.OrderRpc.UpdateChatOrderUse(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.DoChatOrderActionResp{}
	return
}

func (l *DoChatOrderActionLogic) defaultAction(req *types.DoChatOrderActionReq) (resp *types.DoChatOrderActionResp, err error) {
	var in pb.DoChatOrderActionReq
	_ = copier.Copy(&in, req)
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.DoChatOrderActionResp{}
	return
}

func (l *DoChatOrderActionLogic) applyRefund(req *types.DoChatOrderActionReq) (resp *types.DoChatOrderActionResp, err error) {
	var in pb.DoChatOrderActionReq
	_ = copier.Copy(&in, req)
	l.getUserRefundArgs(&in)
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.DoChatOrderActionResp{}
	return
}

// 判断自动退款
func (l *DoChatOrderActionLogic) getUserRefundArgs(req *pb.DoChatOrderActionReq) {
	var order *pb.GetChatOrderDetailResp
	var userStat *pbUser.GetUserStatResp
	err := mr.Finish(func() error {
		var err2 error
		order, err2 = l.svcCtx.OrderRpc.GetChatOrderDetail(l.ctx, &pb.GetChatOrderDetailReq{OrderId: req.OrderId})
		if err2 != nil {
			return err2
		}
		return nil
	},
		func() error {
			var err2 error
			userStat, err2 = l.svcCtx.UsercenterRpc.GetUserStat(l.ctx, &pbUser.GetUserStatReq{
				Uid: req.OperatorUid,
			})
			if err2 != nil {
				return err2
			}
			return nil
		})
	if err != nil {
		return
	}

	// 校验状态
	if !tool.IsInt64ArrayExist(order.Order.OrderState, orderkey.CanApplyRefundOrderState) {
		return
	}
	// 判断是否符合自动同意退款
	createTime, err := time.ParseInLocation(db.DateTimeFormat, order.Order.CreateTime, time.Local)
	if err != nil {
		return
	}
	if createTime.Add(userkey.NoCondRefundTimeLimitMin*time.Minute).After(time.Now()) && userStat.NoCondRefundCnt > 1 { // 是否符合无条件退款
		var rs *pbUser.UpdateNoCondRefundCntResp
		rs, err = l.svcCtx.UsercenterRpc.UpdateNoCondRefundCnt(l.ctx, &pbUser.UpdateNoCondRefundCntReq{Uid: req.OperatorUid})
		if err != nil {
			return
		}
		if rs.State == db.Enable {
			req.Action = orderkey.ChatOrderStateAutoAgreeRefund6
		}
	}
	return
}
