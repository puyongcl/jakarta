package zhihu

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"jakarta/common/xerr"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RestApiClient struct {
	Client *resty.Client
}

var restClient *RestApiClient

func InitZhihuApiClient() *RestApiClient {
	if restClient != nil {
		return restClient
	}
	restClient = &RestApiClient{
		Client: resty.NewWithClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	}
	return restClient
}

const eventHold = "__EVENTTYPE__"
const valueHold = "__VALUE__"
const stampHold = "__TIMESTAMP__"
const cstampHold = "__TS__"
const extra = `&source=zhihu&`

func (c *RestApiClient) UploadEvent(cb, event, value, stamp string) error {
	if cb == "" || event == "" {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	cb, err := url.QueryUnescape(cb)
	if err != nil {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("UploadEvent QueryUnescape %+v", err))
	}

	cb = strings.Replace(cb, extra, "", -1)

	cb = strings.Replace(cb, eventHold, event, 1)
	cb = strings.Replace(cb, valueHold, value, 1)
	cb = strings.Replace(cb, stampHold, stamp, 1)
	cb = strings.Replace(cb, cstampHold, fmt.Sprintf("%d", time.Now().Unix()), 1)
	rsp, err := c.Client.R().Get(cb)
	if err != nil {
		return err
	}
	if rsp.IsError() {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, rsp.Status())
	}
	return nil
}
