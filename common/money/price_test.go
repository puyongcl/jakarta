package money

import (
	"fmt"
	"testing"
)

func TestYuan(t *testing.T) {
	var in int64 = 192342
	amount := fmt.Sprintf("%.2f", float64(in)/MoneyYuan)
	fmt.Println(float64(in) / MoneyYuan)
	fmt.Println(amount)
	fmt.Println(GetYuan(in))
}

func TestRoundYuan(t *testing.T) {
	fmt.Println(RoundYuan(311))
	fmt.Println(RoundYuan(200))
	fmt.Println(RoundYuan(0))
	fmt.Println(RoundYuan(90))
}

func TestCalcShareAmount(t *testing.T) {
	gotPlatformAmount, gotAmount := CalcShareAmount(5, 1, 0, &CurrentShareArg{
		//BaseAmount:           10000,
		TaxAmount: 2000,
		//NightAddAmount:       1000,
		//SaveAmount:           1000,
		ActualAmount:         12000,
		ShareRateStep1Star5:  900,
		ShareRateStep1Star3:  800,
		ShareRateStep1Star1:  100,
		ShareAmountStep1Unit: 1,
		ShareRateStep2Star5:  800,
		ShareRateStep2Star3:  600,
		ShareRateStep2Star1:  200,
	})
	fmt.Println(gotPlatformAmount, gotAmount)
	gotPlatformAmount, gotAmount = CalcShareAmount(5, 1, 1, &CurrentShareArg{
		//BaseAmount:           10000,
		TaxAmount: 2000,
		//NightAddAmount:       1000,
		//SaveAmount:           1000,
		ActualAmount:         12000,
		ShareRateStep1Star5:  900,
		ShareRateStep1Star3:  800,
		ShareRateStep1Star1:  100,
		ShareAmountStep1Unit: 1,
		ShareRateStep2Star5:  800,
		ShareRateStep2Star3:  600,
		ShareRateStep2Star1:  200,
	})
	fmt.Println(gotPlatformAmount, gotAmount)
}

func TestCalcShareAmount2(t *testing.T) {
	gotPlatformAmount, gotAmount := CalcShareAmount(1, 1, 0, &CurrentShareArg{
		//BaseAmount:           10000,
		TaxAmount: 238,
		//NightAddAmount:       1000,
		//SaveAmount:           1000,
		ActualAmount:         3900,
		ShareRateStep1Star5:  400,
		ShareRateStep1Star3:  500,
		ShareRateStep1Star1:  800,
		ShareAmountStep1Unit: 1,
		ShareRateStep2Star5:  200,
		ShareRateStep2Star3:  400,
		ShareRateStep2Star1:  700,
	})
	fmt.Println(gotPlatformAmount, gotAmount)
	gotPlatformAmount, gotAmount = CalcShareAmount(5, 1, 1, &CurrentShareArg{
		//BaseAmount:           10000,
		TaxAmount: 119,
		//NightAddAmount:       1000,
		//SaveAmount:           1000,
		ActualAmount:         2000,
		ShareRateStep1Star5:  400,
		ShareRateStep1Star3:  500,
		ShareRateStep1Star1:  800,
		ShareAmountStep1Unit: 1,
		ShareRateStep2Star5:  200,
		ShareRateStep2Star3:  400,
		ShareRateStep2Star1:  700,
	})
	fmt.Println(gotPlatformAmount, gotAmount)
}

func TestGetCurrentListenerChatPrice(t *testing.T) {
	req := CurrentListenerChatPriceArg{
		TextChatPrice:          3000,
		VoiceChatPrice:         3000,
		TaxRate:                66,
		NightAddPriceRate:      200,
		NightAddPriceHourStart: 13,
		NightAddPriceHourEnd:   18,
		NewUserDiscount:        500,
		ChatUnitMinute:         15,
	}
	got := GetCurrentListenerChatPrice(&req)
	fmt.Println(got)
}
