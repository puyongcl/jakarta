package tim

import (
	"jakarta/common/xerr"
	"strconv"
)

const accountDelUrl = "v4/im_open_login_svc/account_delete"

type DelAccountItem struct {
	UserID string `json:"UserID"`
}
type AccountDelReq struct {
	DeleteItem []*DelAccountItem `json:"DeleteItem"`
}

type AccountDelResultItem struct {
	ResultCode int    `json:"ResultCode"`
	ResultInfo string `json:"ResultInfo"`
	UserID     string `json:"UserID"`
}

type AccountDelResp struct {
	ActionStatus string                  `json:"ActionStatus"`
	ErrorCode    uint32                  `json:"ErrorCode"`
	ErrorInfo    string                  `json:"ErrorInfo"`
	ResultItem   []*AccountDelResultItem `json:"ResultItem,optional"`
}

func (ra *RestApiClient) AccountDel(uid int64) (err error) {
	req := AccountDelReq{
		DeleteItem: []*DelAccountItem{
			{UserID: strconv.FormatInt(uid, 10)},
		},
	}
	var rsp AccountDelResp
	resp, err := ra.Client.R().SetBody(&req).SetResult(&rsp).Post(ra.GetTimRestApiUrl(accountDelUrl))
	if err != nil {
		return
	}
	if resp.IsError() {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, resp.Status())
	}
	if rsp.ActionStatus == "OK" {
		return
	}
	return xerr.NewGrpcErrCodeMsg(rsp.ErrorCode, rsp.ErrorInfo)
}
