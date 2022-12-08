package svc

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"jakarta/common/config"
)

func NewWxpayNotifyHandler(c *config.WxPayCallbackConf) *notify.Handler {
	ctx := context.Background()
	mchPrivateKey, err := utils.LoadPrivateKey(c.PrivateKey)
	if err != nil {
		panic(err)
	}
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, c.SerialNo, c.MchId, c.APIv3Key)
	if err != nil {
		panic(err)
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(c.MchId)
	// 3. 使用证书访问器初始化 `notify.Handler`
	return notify.NewNotifyHandler(c.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
}
