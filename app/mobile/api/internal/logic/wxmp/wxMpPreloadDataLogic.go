package wxmp

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/logic/chatorder"
	"jakarta/app/mobile/api/internal/logic/listener"
	"jakarta/app/mobile/api/internal/logic/user"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"
)

type WxMpPreloadDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxMpPreloadDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxMpPreloadDataLogic {
	return &WxMpPreloadDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxMpPreloadDataLogic) WxMpPreloadData(req *types.GetWxMpPreloadDataReq) (resp *types.GetWxMpPreloadDataResp, err error) {
	req.AppVer = -l.svcCtx.Config.AppVerConf.MinAppVer

	return l.doWxMpLogin(req)
}

func (l *WxMpPreloadDataLogic) WxMpLogin(req *types.GetWxMpPreloadDataReq) (resp *types.GetWxMpPreloadDataResp, err error) {

	return l.doWxMpLogin(req)
}

func (l *WxMpPreloadDataLogic) doWxMpLogin(req *types.GetWxMpPreloadDataReq) (resp *types.GetWxMpPreloadDataResp, err error) {
	resp = &types.GetWxMpPreloadDataResp{
		LoginData:                nil,
		CloudConfigData:          nil,
		ListenerDraftData:        nil,
		Listener:                 nil,
		BusinessConfigData:       nil,
		ChatOrderPriceConfigData: nil,
	}
	// 登陆
	if req.Code == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "code is empty")
	}
	l1 := user.NewWxMiniAuthLogic(l.ctx, l.svcCtx)
	req1 := types.WXMiniAuthReq{
		Code:          req.Code,
		IV:            "",
		EncryptedData: "",
		Query:         req.Query,
		AppVer:        req.AppVer,
	}
	var rsp1 *types.WXMiniAuthResp
	rsp1, err = l1.WxMiniAuth(&req1)
	if err != nil {
		return nil, err
	}
	logx.WithContext(l.ctx).Infof("WxMpPreloadDataLogic user wx LOGIN USERTYPE:%d resp:%+v, user:%+v", rsp1.UserType, rsp1, rsp1.User)
	resp.LoginData = rsp1

	//
	var rsp2 *types.GetCloudConfigResp
	var rsp3 *types.GetListenerOwnInfoResp
	var rsp4 *types.RecommendListenerResp
	var rsp5 *types.GetDefineBusinessConfigResp
	var rsp6 *types.GetBusinessChatPricingPlanResp
	var rsp7 *types.GetListenerHomePageDashboardResp
	var rsp8 *types.GetListenerSeeChatOrderListResp

	err = mr.Finish(func() error {
		var err2 error
		req2 := types.GetCloudConfigReq{Uid: rsp1.User.Uid}
		l2 := user.NewGetCloudConfigLogic(l.ctx, l.svcCtx)
		rsp2, err2 = l2.GetCloudConfig(&req2)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		req3 := types.GetListenerOwnInfoReq{ListenerUid: rsp1.User.Uid}
		l2 := listener.NewGetListenerOwnProfileLogic(l.ctx, l.svcCtx)
		rsp3, err2 = l2.GetListenerOwnProfile2(&req3, false)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		req4 := types.RecommendListenerReq{
			Uid:         rsp1.User.Uid,
			PageNo:      1,
			PageSize:    10,
			Specialties: 0,
			ChatType:    0,
			Gender:      0,
			Age:         0,
			SortOrder:   0,
		}
		l2 := listener.NewGetRecommendListenerListLogic(l.ctx, l.svcCtx)
		rsp4, err2 = l2.DoGetRecommendListenerList(&req4, rsp1.OpenId)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		req5 := types.GetDefineBusinessConfigReq{Uid: rsp1.User.Uid}
		l2 := listener.NewGetDefineBusinessConfigLogic(l.ctx, l.svcCtx)
		rsp5, err2 = l2.GetBannerAndFilterConfigReq(&req5)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		req6 := types.GetBusinessChatPricingPlanReq{Uid: rsp1.User.Uid}
		l2 := chatorder.NewGetBusinessChatPriceLogic(l.ctx, l.svcCtx)
		rsp6, err2 = l2.GetBusinessChatPrice(&req6)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		if rsp1.UserType != userkey.UserTypeListener {
			return nil
		}
		req7 := types.GetListenerHomePageDashboardReq{ListenerUid: rsp1.User.Uid}
		l2 := listener.NewGetListenerHomePageDashboardLogic(l.ctx, l.svcCtx)
		rsp7, err2 = l2.GetListenerHomePageDashboard2(&req7, false)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var err2 error
		if rsp1.UserType != userkey.UserTypeListener {
			return nil
		}
		req8 := types.GetListenerSeeChatOrderListReq{
			ListenerUid: rsp1.User.Uid,
			ListType:    orderkey.OrderListTypeNeedProcess,
			PageNo:      1,
			PageSize:    3,
		}
		l2 := chatorder.NewListenerChatOrderListLogic(l.ctx, l.svcCtx)
		rsp8, err2 = l2.ListenerChatOrderList(&req8)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 登陆时推荐XXX、发送问候语
	if rsp1.UserType == userkey.UserTypeNormalUser {
		resp.RecListener = l.recListener(rsp1.User.Uid, rsp1.IsNewUser, rsp1.User.CreateTime, rsp1.OpenId, rsp6.Config)
	}

	resp.LoginData = rsp1
	resp.CloudConfigData = rsp2
	if rsp3 != nil {
		resp.ListenerDraftData = rsp3.Info
	}
	if rsp4 != nil {
		resp.Listener = rsp4.Listener
	}
	resp.BusinessConfigData = rsp5
	resp.ChatOrderPriceConfigData = rsp6.Config
	resp.ListenerHomeData = rsp7
	if rsp8 != nil {
		resp.ListenerHomeOrderListData = rsp8.List
	}

	return
}

func (l *WxMpPreloadDataLogic) recListener(uid, isNewUser int64, createTime string, authKey string, bp *types.BusinessChatPricingPlan) *types.UserSeeRecommendListenerProfile {
	var err error
	var regDays int64
	if isNewUser != db.Enable {
		var regTime time.Time
		regTime, err = time.ParseInLocation(db.DateTimeFormat, createTime, time.Local)
		if err != nil {
			return nil
		}
		regDays = int64(time.Now().Sub(regTime).Hours())
		regDays = tool.DivideInt64(regDays, 24)
	}

	in := pbListener.RecListenerWhenUserLoginReq{
		Uid:       uid,
		IsNewUser: isNewUser,
		RegDays:   regDays,
		OrderCnt:  bp.OrderCnt,
	}
	var rsp *pbListener.RecListenerWhenUserLoginResp
	rsp, err = l.svcCtx.ListenerRpc.RecListenerWhenUserLogin(l.ctx, &in)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("WxMpPreloadDataLogic RecListenerWhenUserLogin err:%+v", err)
		return nil
	}
	if rsp.RecListener == nil {
		return nil
	}
	bpp := pbOrder.BusinessChatPricePlan{}
	_ = copier.Copy(&bpp, bp)

	return listener.AddListenerData(rsp.RecListener, &bpp)
}
