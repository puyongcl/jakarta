package tim

import (
	"jakarta/common/xerr"
	"strconv"
)

type BlockUserReq struct {
	FromAccount string   `json:"From_Account"` // 必填	请求为该 UserID 添加黑名单
	ToAccount   []string `json:"To_Account"`   // 必填	待添加为黑名单的用户 UserID 列表，单次请求的 To_Account 数不得超过1000
}

type BlockUserResp struct {
	ResultItem []struct {
		ToAccount  string `json:"To_Account"`
		ResultCode int    `json:"ResultCode"`
		ResultInfo string `json:"ResultInfo"`
	} `json:"ResultItem"`
	FailAccount  []string `json:"Fail_Account"`
	ActionStatus string   `json:"ActionStatus"`
	ErrorCode    uint32   `json:"ErrorCode"`
	ErrorInfo    string   `json:"ErrorInfo"`
	ErrorDisplay string   `json:"ErrorDisplay"`
}

const addBlacklistUrl = "v4/sns/black_list_add"
const deleteBlacklistUrl = "v4/sns/black_list_delete"

func (ra *RestApiClient) AddBlacklist(fromUid int64, blockUid int64) error {
	return ra.blacklist(fromUid, blockUid, addBlacklistUrl)
}

func (ra *RestApiClient) DeleteBlacklist(fromUid int64, blockUid int64) error {
	return ra.blacklist(fromUid, blockUid, deleteBlacklistUrl)
}

func (ra *RestApiClient) blacklist(fromUid int64, blockUid int64, url string) (err error) {
	req := BlockUserReq{
		FromAccount: strconv.FormatInt(fromUid, 10),
		ToAccount:   []string{strconv.FormatInt(blockUid, 10)},
	}
	var rsp BlockUserResp
	resp, err := ra.Client.R().SetBody(&req).SetResult(&rsp).Post(ra.GetTimRestApiUrl(url))
	if err != nil {
		return
	}
	if resp.IsError() {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, resp.Status())
	}
	if rsp.ActionStatus == "OK" {
		if len(rsp.ResultItem) == 1 {
			if rsp.ResultItem[0].ResultCode != 0 {
				return xerr.NewGrpcErrCodeMsg(rsp.ErrorCode, rsp.ErrorInfo)
			}
		}
		return
	}
	return xerr.NewGrpcErrCodeMsg(rsp.ErrorCode, rsp.ErrorInfo)
}
