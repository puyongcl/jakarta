package tencentcloud

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"jakarta/common/xerr"
	"net/http"
	"net/url"
)

type TencentCosClient struct {
	client *cos.Client
	bucket string
	region string
}

// 上传文件状态
const (
	UploadStateInit    = 1
	UploadStateSuccess = 2
	UploadStateFail    = 3
)

func InitCosClient(sid, skey, bucket, region string) *TencentCosClient {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: sid,
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: skey,
		},
	})

	return &TencentCosClient{
		client: client,
		bucket: bucket,
		region: region,
	}
}

func (c *TencentCosClient) Upload(ctx context.Context, key, localFile string) (err error) {
	var rsp *cos.Response
	_, rsp, err = c.client.Object.Upload(ctx, key, localFile, nil)
	if err != nil {
		err = xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%+v", err))
		return
	}
	if rsp != nil && rsp.StatusCode != http.StatusOK {
		err = xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%s", rsp.Status))
	}
	return
}
