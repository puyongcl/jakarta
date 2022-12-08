package tim

import (
	"fmt"
	"jakarta/common/key/tencentcloudkey"
	"testing"
)

const testSdkAppId = 144778745455445
const testAdminId = "administrator"
const testAdminSig = "eJwtzMsOgjAUBNB-3211313131333-2-oHlityvXdnX6lI-d8nh*5-6BvaLLvMjcZh-kEGFUZGthF*T7A45aNNs_"

func TestInitTimRestApiClient(t *testing.T) {
	tc := InitTimRestApiClient(testSdkAppId, testAdminId, testAdminSig)
	err := tc.AccountImport(12, "", "")
	fmt.Println(err)
}

func TestRestApiClient_UpdateProfile(t *testing.T) {
	tc := InitTimRestApiClient(testSdkAppId, testAdminId, testAdminSig)
	err := tc.UpdateProfile(100002, "TestName", tencentcloudkey.CDNBasePath+"avatar/FskoR7AjXPHLf8c1bf4aa811f1f81d2b34dac25f61f3.jpeg")
	if err != nil {
		fmt.Println(err)
	}
}

func TestRestApiClient_AccountDel(t *testing.T) {
	tc := InitTimRestApiClient(testSdkAppId, testAdminId, testAdminSig)
	err := tc.AccountDel(100025)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRestApiClient_SendMsg(t *testing.T) {
	tc := InitTimRestApiClient(testSdkAppId, testAdminId, testAdminSig)
	_, err := tc.SendTextMsg(200047, 200031, "你好3", TimMsgSyncFromNo)
	if err != nil {
		fmt.Println(err)
	}
}
