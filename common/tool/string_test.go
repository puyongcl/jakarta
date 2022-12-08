package tool

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCutText(t *testing.T) {
	fmt.Println(CutText("好好学习是需要的，我们是要一直提升自己，自我突破的，包括我现在也每天都在学习。加油！如果再有什么想聊的就上来跟我吐槽也可以[龇牙]。", 20, "..."))
}

func TestRemoveSpace(t *testing.T) {
	txt := "622208 4000009937627"
	fmt.Println(strings.Replace(txt, " ", "", -1))
}

func TestSign(t *testing.T) {
	type HFBFCashCallbackData struct {
		Type            int64    `json:"type"`             // 1 支付结果回调 2 任务作废回调
		WorkNumber      string   `json:"work_number"`      // 任务编号
		CompanyId       int64    `json:"company_id"`       //
		Status          string   `json:"status"`           // -1 作废
		TovoidTime      string   `json:"tovoid_time"`      // 作废时间
		CustomerNumbers []string `json:"customer_numbers"` // 此任务包含的所有自定义任务单号
		UserId          string   `json:"user_id"`          // 用户id
		Number          string   `json:"number"`           // 打款流水号
		PayStatus       string   `json:"pay_status"`       // 打款状态 1 未结算 2 待结算 3 结算中 4 已结算 5 结算失败
		CustomNumber    string   `json:"custom_number"`    // 自定义流水号
		Msg             string   `json:"msg"`              // 支付失败原因
		PayTime         string   `json:"pay_time"`         // 支付时间 格式2020-10-10 12:00:00
		Sign            string   `json:"sign"`             // 签名 验证参数 值为 timeStamp 加 params ⾥所有字段值 拼起来再加appsecret求 MD5，忽略空值，转成⼤写。字 段按字⺟表由⼩到⼤排序。右边示例展示了密码为 123abcd 的 Sign 的⽣成过程。请注意 Bool 类型有⼤⼩写：True/False,Array 类型直接拼接，如[1,2,30]=1230
	}

	data := HFBFCashCallbackData{
		Type:            1,
		WorkNumber:      "22",
		CompanyId:       2,
		Status:          "332",
		TovoidTime:      "2121",
		CustomerNumbers: []string{"adadad", "ssdadad"},
		UserId:          "1dads",
		Number:          "dadsa",
		PayStatus:       "dadsa",
		CustomNumber:    "dadsada",
		Msg:             "dadsa",
		PayTime:         "dasd",
		Sign:            "das",
	}
	var cc map[string]interface{}
	b, _ := json.Marshal(&data)
	_ = json.Unmarshal(b, &cc)
	fmt.Println(cc)
}

func TestCountString(t *testing.T) {
	tt := "043YHD0w3vZZiZ2v2O1w38rgs81YHD06"
	ll := utf8.RuneCountInString(tt)
	fmt.Println(len(tt))

	fmt.Println(ll)
}
