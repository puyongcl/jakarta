package svc

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"jakarta/common/config"
)

func NewWxPayClientV3(c *config.WxPayConf) *core.Client {
	mchPrivateKey, err := utils.LoadPrivateKey(c.PrivateKey)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// Initialize the client with the merchant's private key, etc., and make it have the ability to automatically obtain WeChat payment platform certificates at regular intervals
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.MchId, c.SerialNo, mchPrivateKey, c.APIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		panic(err)
	}

	return client
}
