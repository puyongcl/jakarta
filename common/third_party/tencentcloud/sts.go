package tencentcloud

import (
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"time"
)

type TencentStsClient struct {
	client *sts.Client
	appid  string
	bucket string
	region string
	expire int64
}

func InitStsClient(sid, skey, appid, bucket, region string, expire int64) *TencentStsClient {
	client := &TencentStsClient{appid: appid, bucket: bucket, region: region, expire: expire}
	c := sts.NewClient(
		sid,
		skey,
		nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	client.client = c
	return client
}

func (c *TencentStsClient) GetTempSecret() (*sts.CredentialResult, error) {
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          c.region,
		Policy: &sts.CredentialPolicy{
			Version: "2.0",
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PutObject",
						"name/cos:PostObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + c.region + ":uid/" + c.appid + ":" + c.bucket + "/*",
					},
				},
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:GetObject",
						"name/cos:HeadObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + c.region + ":uid/" + c.appid + ":" + c.bucket + "/*",
					},
				},
			},
		},
	}

	// case 1 请求临时密钥
	res, err := c.client.GetCredential(opt)
	if err != nil {
		return nil, err
	}
	return res, err
}
