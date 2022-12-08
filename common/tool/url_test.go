package tool

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestUrlQueryDecode(t *testing.T) {
	q := `source=zhihu\u0026cb=https%3A%2F%2Fsugar.zhihu.com%2Fplutus_adreaper_callback%3Fsi%3D5eae44a2-ba85-4029-83ac-1e16a08c3d4e%26os%3D0%26zid%3D105%26zaid%3D2543273%26zcid%3D2511344%26cid%3D2511344%26event%3D__EVENTTYPE__%26value%3D__EVENTVALUE__%26ts%3D__TIMESTAMP__%26cts%3D__TS__%26mh%3Ded13c8e3bf465920390b0842d5ec2d3d`
	//p := `source=zhihu\u0026cb=https://sugar.zhihu.com/plutus_adreaper_callback?si=5eae44a2-ba85-4029-83ac-1e16a08c3d4e&os=0&zid=105&zaid=2543273&zcid=2511344&cid=2511344&event=__EVENTTYPE__&value=__EVENTVALUE__&ts=__TIMESTAMP__&cts=__TS__&mh=ed13c8e3bf465920390b0842d5ec2d3d`
	//q = strings.Replace(q, "\\u0026", "", -1)
	got, got2, err := UrlQueryDecode(q)
	fmt.Println(got, got2, err)

	q1 := `uctrackid=czo2NzA5OTQxMjIxODA5MzQ0OTI4O2M6Nzc0MDg2MzE7ZDpkbXBfNTI4MzE1MDg4OTE0MTA3ODE3NztwOmhj\u0026compId=3E235B31-C404-410C-8E17-485CB49631B2\u0026source=uc\u0026userId=210980664\u0026sid=6709941221809344928\u0026`
	got, got2, err = UrlQueryDecode(q1)
	fmt.Println(got, got2, err)

	q3 := `fid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26source%3Dbaidu%26pageId%3D114293509%26bd_vid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26`
	got, got2, err = UrlQueryDecode(q3)
	fmt.Println(got, got2, err)
}

func TestUrlQueryDecode2(t *testing.T) {
	q3 := `fid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26source%3Dbaidu%26pageId%3D114293509%26bd_vid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26`
	got, got2, err := UrlQueryDecode(q3)
	fmt.Println(got, got2, err)
}

func TestUrlDecode(t *testing.T) {
	q := `timeStamp=1663567813\u0026params%5Btype%5D=1\u0026params%5Bcompany_id%5D=92\u0026params%5Buser_id%5D=93f4a9b768ddb1c1167e951cf0f0151a\u0026params%5Bwork_number%5D=K20220919120038UEKHVH\u0026params%5Bnumber%5D=Z20220919120038W0ZT0W1\u0026params%5Bpay_status%5D=4\u0026params%5Bcustom_number%5D=CW20220916165137103e900000c\u0026params%5Bmsg%5D=-\u0026params%5Bpay_time%5D=2022-09-19+12%3A03%3A03\u0026sign=0667E5CE35FB66F31A86FC74FA75900F`
	ul, err := url.QueryUnescape(q)
	if err != nil {
		fmt.Println(err)
		return
	}
	ul = strings.Replace(ul, "\\u0026", "&", -1)
	//fmt.Println(ul)
	// 解析
	//"timeStamp=1663567813&params[type]=1&params[company_id]=92&params[user_id]=93f4a9b768ddb1c1167e951cf0f0151a&params[work_number]=K20220919120038UEKHVH&params[number]=Z20220919120038W0ZT0W1&params[pay_status]=4&params[custom_number]=CW20220916165137103e900000c&params[msg]=-&params[pay_time]=2022-09-19 12:03:03&sign=0667E5CE35FB66F31A86FC74FA75900F"
	arr := strings.Split(ul, "&")
	ma := make(map[string]string)
	for idx := 0; idx < len(arr); idx++ {
		as := strings.Split(arr[idx], "=")
		if len(as) == 2 {
			k := strings.Replace(as[0], "params[", "", -1)
			k = strings.Replace(k, "]", "", -1)
			ma[k] = as[1]
		}
	}
	fmt.Println(ma)
}
