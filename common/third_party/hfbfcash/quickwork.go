package hfbfcash

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"jakarta/common/xerr"
	"net/http"
	"strconv"
	"time"
)

// 快捷打款
type QuickWorkDetail struct {
	Name         string `json:"name"`
	Idcard       string `json:"idcard"`
	Bank         string `json:"bank"`
	BankCard     string `json:"bank_card"`
	Phone        string `json:"phone"`
	Price        string `json:"price"`
	CustomNumber string `json:"custom_number"`
}

type QuickWorkReq struct {
	Industry   string            `json:"industry"`
	WorkName   string            `json:"work_name"`
	TotalPrice string            `json:"total_price"`
	Detail     []QuickWorkDetail `json:"detail"`
}

type QuickWorkRsp struct {
	Code    int64         `json:"code"`
	Message string        `json:"message"`
	Data    QuickWorkData `json:"data"`
}

type QuickWorkData struct {
	WorkNumber string `json:"work_number"`
}

const quickWorkAPi = ""

// 快捷银行卡打款
func (r *RestHfbfCashClient) QuickWork(req *QuickWorkReq) (resp *QuickWorkRsp, err error) {
	resp = &QuickWorkRsp{}
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	encrypted, err := aesEncrypt(buf, []byte(r.AppSecret), []byte(r.AppId)) //加密
	if err != nil {
		panic(err)
	}
	m := map[string]string{
		"appId":     r.AppId,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"data":      base64.StdEncoding.EncodeToString(encrypted),
	}

	rs, err := r.Client.R().SetFormData(m).SetResult(resp).Post(baseHost + quickWorkAPi)
	if err != nil {
		return
	}
	if rs.IsError() {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, rs.Status())
	}
	if resp.Code != http.StatusOK {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, fmt.Sprintf("%d-%s", resp.Code, resp.Message))
	}
	return
}
