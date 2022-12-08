package logic

import (
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/paykey"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建微信支付预处理订单
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {
	// 查询当前订单是否有支付流水
	cnt, err := l.svcCtx.ThirdPaymentFlowModel.FindCount(l.ctx, in.OrderId, paykey.NotNeedPayState)
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.PaymentErrorFlowAlreadyExist, "已经支付")
	}
	//
	data := new(paymentPgModel.ThirdPaymentFlow)
	data.FlowNo = uniqueid.GenSn(uniqueid.SnPrefixThirdPaymentFlowNo)
	data.Uid = in.Uid
	data.PayMode = in.PayModel
	data.PayAmount = in.PayAmount
	data.OrderId = in.OrderId
	data.OrderType = in.OrderType
	data.PayStatus = paykey.ThirdPaymentPayTradeStateWait

	//
	var wxrsp *jsapi.PrepayWithRequestPaymentResponse
	err = l.svcCtx.ThirdPaymentFlowModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var err2 error
		_, err2 = l.svcCtx.ThirdPaymentFlowModel.InsertTrans(ctx, session, data)
		if err2 != nil {
			return err2
		}

		// create wechat pay pre pay order
		jsApiSvc := jsapi.JsapiApiService{Client: l.svcCtx.WxPayClient}
		// Get the prepay_id, as well as the parameters and signatures needed to invoke the payment
		timeExpire := time.Now().Add(paykey.PaymentExpireTimeMinute * time.Minute)
		jsapiReq := jsapi.PrepayRequest{
			TimeExpire:  &timeExpire,
			Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
			Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
			Description: core.String(in.OrderId),
			OutTradeNo:  core.String(data.FlowNo),
			Attach:      core.String(in.Description),
			NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),
			Amount: &jsapi.Amount{
				Total:    core.Int64(in.PayAmount),
				Currency: core.String("CNY"),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(in.OpenId),
			},
		}
		wxrsp, _, err2 = jsApiSvc.PrepayWithRequestPayment(ctx, jsapiReq)
		if err2 != nil {
			return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%+v", err2))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := &pb.CreatePaymentResp{
		Appid:     *wxrsp.Appid,
		NonceStr:  *wxrsp.NonceStr,
		PaySign:   *wxrsp.PaySign,
		Package:   *wxrsp.Package,
		Timestamp: *wxrsp.TimeStamp,
		SignType:  *wxrsp.SignType,
	}
	return resp, nil
}
