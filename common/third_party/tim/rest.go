package tim

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"time"
)

const (
	timBaseUrl = "https://console.tim.qq.com/"

	restApiArgUrl = "?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=json"
)

type RestApiClient struct {
	SdkAppid       int64
	AdminId        string
	AdminSignature string
	randInt        *rand.Rand
	Client         *resty.Client
}

type CommonResp struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    uint32 `json:"ErrorCode"`
}

var timClient *RestApiClient

func InitTimRestApiClient(sdkappid int64, adminId string, adminSig string) *RestApiClient {
	if timClient != nil {
		return timClient
	}
	timClient = &RestApiClient{
		SdkAppid:       sdkappid,
		AdminId:        adminId,
		AdminSignature: adminSig,
		randInt:        rand.New(rand.NewSource(time.Now().UnixNano())),
		Client: resty.NewWithClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
	}
	return timClient
}

func GetTimClient() *RestApiClient {
	return timClient
}

func (ra *RestApiClient) GetTimRestApiUrl(reqUrl string) string {
	return timBaseUrl + reqUrl + fmt.Sprintf(restApiArgUrl, ra.SdkAppid, ra.AdminId, ra.AdminSignature, ra.randInt.Uint32())
}
