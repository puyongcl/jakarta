package hfbfcash

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"jakarta/common/xerr"
	"net/http"
	"strconv"
	"time"
)

const syncContractAPI = ""

type SyncContractReq struct {
	Name   string            `json:"name"`   // 姓名
	Phone  string            `json:"phone"`  // 电话号码
	Idcard string            `json:"idcard"` // 身份证号
	Remark string            `json:"remark"`
	Files  *SyncContractFile `json:"files"`
}

type SyncContractFile struct {
	File1 string `json:"file1"`
	File2 string `json:"file2"`
	File3 string `json:"file3"`
}

type SyncContractData struct {
}
type SyncContractResp struct {
	Code    int64                 `json:"code"`
	Message string                `json:"message"`
	Data    CustomerWorkQueryData `json:"data"` // 用户合同地址
}

// 同步合同
func (r *RestHfbfCashClient) SyncContract(req *SyncContractReq) (err error) {
	resp := &SyncContractResp{}
	var buf []byte
	buf, err = json.Marshal(req)
	if err != nil {
		return err
	}
	var encrypted []byte
	encrypted, err = aesEncrypt(buf, []byte(r.AppSecret), []byte(r.AppId)) //加密
	if err != nil {
		panic(err)
	}
	m := map[string]string{
		"appId":     r.AppId,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"data":      base64.StdEncoding.EncodeToString(encrypted),
	}

	var rs *resty.Response
	rs, err = r.Client.R().SetFormData(m).SetResult(resp).Post(baseHost + syncContractAPI)
	if err != nil {
		err = xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%+v", err))
		return
	}
	if rs.IsError() {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, rs.Status())
	}
	if resp.Code != http.StatusOK {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%d-%s", resp.Code, resp.Message))
	}
	return
}
