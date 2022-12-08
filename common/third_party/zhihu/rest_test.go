package zhihu

import (
	"fmt"
	"testing"
	"time"
)

func TestRestApiClient_UploadEvent(t *testing.T) {
	zrc := InitZhihuApiClient()
	cb := `https://sugar.zhihu.com/plutus_adreaper_callback?si=5eae44a2-ba85-4029-83ac-1e16a08c3d4e&os=0&zid=105&zaid=2543273&zcid=2511344&cid=2511344&event=__EVENTTYPE__&value=__EVENTVALUE__&ts=__TIMESTAMP__&cts=__TS__&mh=ed13c8e3bf465920390b0842d5ec2d3d`
	//cb2 := `https%3A%2F%2Fsugar.zhihu.com%2Fplutus_adreaper_callback%3Fsi%3D6810a311-e91a-4dd6-babc-4d2b594dab67%26os%3D1%26zid%3D105%26zaid%3D2543273%26zcid%3D2513687%26cid%3D2513687%26event%3D__EVENTTYPE__%26value%3D__EVENTVALUE__%26ts%3D__TIMESTAMP__%26cts%3D__TS__%26mh%3D1b6c08eebabc6541ec405c10dbb1ab46\\u0026source=zhihu\\u0026`
	err := zrc.UploadEvent(cb, "api_view", "", fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
