package svc

import (
	_ "github.com/lib/pq"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/payment/rpc/internal/config"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/db"
	"jakarta/common/third_party/hfbfcash"
)

type ServiceContext struct {
	Config                config.Config
	ThirdPaymentFlowModel paymentPgModel.ThirdPaymentFlowModel
	ThirdRefundFlowModel  paymentPgModel.ThirdRefundFlowModel
	ThirdCashFlowModel    paymentPgModel.ThirdCashFlowModel

	WxPayClient    *core.Client
	HfbfCashClient *hfbfcash.RestHfbfCashClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	return &ServiceContext{
		Config:                c,
		ThirdPaymentFlowModel: paymentPgModel.NewThirdPaymentFlowModel(sqlConn),
		ThirdRefundFlowModel:  paymentPgModel.NewThirdRefundFlowModel(sqlConn),
		ThirdCashFlowModel:    paymentPgModel.NewThirdCashFlowModel(sqlConn),

		WxPayClient:    NewWxPayClientV3(&c.WxPayConf),
		HfbfCashClient: hfbfcash.InitCqCashClient(c.HfbfCashConf.AppId, c.HfbfCashConf.AppSecret),
	}
}
