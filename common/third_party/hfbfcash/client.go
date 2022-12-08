package hfbfcash

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

//
type RestHfbfCashClient struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
	Client    *resty.Client
}

//const baseHost = ":9000"

const callback = "https://third.domain.com:8888/payment/v1/third/hfbfcallback"

const baseHost = ":8000"

var ErrCodeDesc = map[int64]string{
	1001: "参数错误",
	1002: "手机号错误",
	1003: "身份证号码错误",
	1005: "appId错误",
	1006: "任务条数超出限制",
	1007: "银行编号错误",
	1008: "实名认证失败",
	500:  "服务器错误",
	1009: "手机号已存在",
	1010: "IP不在白名单",
}

//const Industry = "IT/互联网"
const Industry = "咨询服务"
const SuccessCode = 200

const DateTimeFormat = "2006-01-02 15:04:05"

var cqClient *RestHfbfCashClient

func InitCqCashClient(appId, appSecret string) *RestHfbfCashClient {
	if cqClient != nil {
		return cqClient
	}
	cqClient = &RestHfbfCashClient{
		AppId:     appId,
		AppSecret: appSecret,
		Client: resty.NewWithClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	}

	return cqClient
}
