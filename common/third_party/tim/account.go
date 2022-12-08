package tim

import (
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"jakarta/common/xerr"
	"strconv"
)

func GenUserSig(userid int64, sdkappid int64, imkey string, expire int) (string, error) {
	uid := strconv.FormatInt(userid, 10)
	return tencentyun.GenUserSig(int(sdkappid), imkey, uid, expire)
}

const accountImportUrl = "v4/im_open_login_svc/account_import"

type AccountImportReq struct {
	UserID  string `json:"UserID"`
	Nick    string `json:"Nick"`
	FaceUrl string `json:"FaceUrl"`
}

func (ra *RestApiClient) AccountImport(uid int64, avatar, nickName string) (err error) {
	req := AccountImportReq{
		UserID:  strconv.FormatInt(uid, 10),
		Nick:    nickName,
		FaceUrl: avatar,
	}
	var rsp CommonResp
	resp, err := ra.Client.R().SetBody(&req).SetResult(&rsp).Post(ra.GetTimRestApiUrl(accountImportUrl))
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
