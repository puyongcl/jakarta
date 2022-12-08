package user

import (
	"context"
	"fmt"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserWxPhoneNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserWxPhoneNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserWxPhoneNumberLogic {
	return &GetUserWxPhoneNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserWxPhoneNumberLogic) GetUserWxPhoneNumber(req *types.GetUserWxPhoneNumerReq) (resp *types.GetUserWxPhoneNumerResp, err error) {
	if req.Code == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "code is empty")
	}

	getPhoneNumRsp, err := l.svcCtx.WxMini.GetAuth().GetPhoneNumberContext(l.ctx, req.Code)
	if err != nil || getPhoneNumRsp.ErrCode != 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.WXminiAuthFail, fmt.Sprintf("获取手机号失败 %+v", err))
	}

	resp = &types.GetUserWxPhoneNumerResp{PhoneNumber: getPhoneNumRsp.PhoneInfo.PhoneNumber, PurePhoneNumber: getPhoneNumRsp.PhoneInfo.PurePhoneNumber, CountryCode: getPhoneNumRsp.PhoneInfo.CountryCode}
	return
}
