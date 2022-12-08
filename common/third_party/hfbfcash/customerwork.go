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

const customerWorkAPI = ""

type CustomerWorkQueryReq struct {
	CustomNumber string `json:"custom_number"`
}

type CustomerWorkQueryData struct {
	WorkNumber   string `json:"work_number"`   // 任务编号
	Number       string `json:"number"`        // 打款流水号
	PayStatus    int64  `json:"pay_status"`    // 打款状态 1未结算 2待结算 3结算中 4已结算(打款成功) 5结算失败
	PayTime      string `json:"pay_time"`      // 打款时间
	CreateTime   string `json:"create_time"`   // 创建时间
	Price        string `json:"price"`         // 打款金额
	Invoice      string `json:"invoice"`       // 发票地址
	UserId       string `json:"user_id"`       // 打款用户id
	CustomNumber string `json:"custom_number"` // 自定义流水号
	Remark       string `json:"remark"`        // 打款备注
	WorkName     string `json:"work_name"`     // 任务名称
	PayMsg       string `json:"pay_nsg"`       // 打款失败原因
	Status       int64  `json:"status"`        // 任务状态 -1作废 1待提交 2待完成 3已完成
}
type CustomerWorkQueryResp struct {
	Code    int64                 `json:"code"`
	Message string                `json:"message"`
	Data    CustomerWorkQueryData `json:"data"`
}

// 快捷银行卡打款
func (r *RestHfbfCashClient) CustomerWorkQuery(req *CustomerWorkQueryReq) (resp *CustomerWorkQueryResp, err error) {
	resp = &CustomerWorkQueryResp{}
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

	rs, err := r.Client.R().SetFormData(m).SetResult(resp).Post(baseHost + customerWorkAPI)
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
