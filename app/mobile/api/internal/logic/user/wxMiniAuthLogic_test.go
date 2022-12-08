package user

import (
	"context"
	"fmt"
	"jakarta/app/mobile/api/internal/types"
	"testing"
)

func Test_parseChannel(t *testing.T) {
	gotChannel, gotCb := parseChannel(context.Background(), &types.WXMiniAuthReq{
		Code:          "",
		IV:            "",
		EncryptedData: "",
		Query:         "https%3A%2F%2Fsugar.zhihu.com%2Fplutus_adreaper_callback%3Fsi%3Dc945ba15-20c7-4208-905a-09b646c5ea51%26os%3D0%26zid%3D105%26zaid%3D2543273%26zcid%3D2513690%26cid%3D2513690%26event%3D__EVENTTYPE__%26value%3D__EVENTVALUE__%26ts%3D__TIMESTAMP__%26cts%3D__TS__%26mh%3Dfcbe950b3451e1bd08fbf0a907fbed2f&",
	})
	fmt.Println(gotChannel, gotCb)
}

func Test_parseChannel2(t *testing.T) {
	gotChannel, gotCb := parseChannel(context.Background(), &types.WXMiniAuthReq{
		Code:          "",
		IV:            "",
		EncryptedData: "",
		//Query:         "fid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26source%3Dbaidu%26pageId%3D114293509%26bd_vid%3DnWmLPj6LPjR1njDdnjn3nWm1nNtkPWmkg17xnH0sg1wxPWfdnHbknHb1n1f%26",
		//Query: "source%3Dbaidu%26",
		Query: "uctrackid%3Dczo0OTQ1NDAyMzMxMDM5Mjg1NDUyO2M6OTc2OTUwNDA5NDtkOmRtcF81NTk1OTQ3Mjg3NDQ3MDI5OTM2O3A6d2w%26compId%3D3E235B31-C404-410C-8E17-485CB49631B2%26source%3Duc%26userId%3D210980664%26sid%3D4945402331039285452%26",
	})
	fmt.Println(gotChannel, gotCb)
}
