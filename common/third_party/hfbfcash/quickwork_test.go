package hfbfcash

import (
	"fmt"
	"testing"
)

func TestRestCqCashClient_quickWork(t *testing.T) {
	req := QuickWorkReq{
		Industry:   Industry,
		WorkName:   "test1",
		TotalPrice: "",
		Detail: []QuickWorkDetail{
			{
				Name:         "t1",
				Idcard:       "614278195812014569",
				BankCard:     "4720680476520721",
				Phone:        "13600421478",
				Price:        "1",
				CustomNumber: "1",
			},
		},
	}

	r := InitCqCashClient("2022072217340493", "DOPMQOZTKJX5WPCU8Y5OMYBEV584QMGH")
	rsp, err := r.QuickWork(&req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func TestRestCqCashClient_syncContract(t *testing.T) {
	req := SyncContractReq{
		Name:   "test",
		Phone:  "1289992929",
		Idcard: "121131313",
		Files: &SyncContractFile{
			File3: "test1",
		},
	}

	r := InitCqCashClient("2022072217340493", "DOPMQOZTKJX5WPCU8Y5OMYBEV584QMGH")
	err := r.SyncContract(&req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
