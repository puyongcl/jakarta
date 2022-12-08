package money

import (
	"fmt"
	"jakarta/common/key/listenerkey"
	"jakarta/common/tool"
	"math"
	"time"
)

// 向上取整
func Round(f float64) int64 {
	return int64(math.Ceil(f))
}

// 除法
func RoundDivideInt64(a int64, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	c := a / b
	return int64(math.Ceil(float64(c)))
}

// 取到元
func RoundYuan(amount int64) int64 {
	if amount <= 0 {
		return 0
	}
	if amount%MoneyYuan > 0 {
		r := amount/MoneyYuan + 1
		return r * MoneyYuan
	}
	return amount
}

// 计算单价价格
func CalculatePrice(price int64, rate int64, div int64) int64 {
	return RoundDivideInt64(price*rate, div)
}

// 计算夜间服务费
func CalculateNightAdd(hour, price int64, rate int64, div int64, startHour, endHour int64) int64 {
	if hour >= startHour && hour <= endHour {
		return CalculatePrice(price, rate, div)
	}
	return 0
}

type CurrentListenerChatPriceArg struct {
	TextChatPrice          int64 `json:"textChatPrice"`          //文字聊天单价格
	VoiceChatPrice         int64 `json:"voiceChatPrice"`         //语音聊天单价格
	TaxRate                int64 `json:"taxRate"`                // 税率(千分之几)
	NightAddPriceRate      int64 `json:"nightAddPriceRate"`      // 夜间加收费率(千分之几)
	NightAddPriceHourStart int64 `json:"nightAddPriceHourStart"` // 夜间服务加价开始时刻
	NightAddPriceHourEnd   int64 `json:"nightAddPriceHourEnd"`   // 夜间加价结束时刻
	NewUserDiscount        int64 `json:"newUserDiscount"`        // 新用户减免比率(千分之几)
	ChatUnitMinute         int64 `json:"chatUnitMinute"`         // 单位时长的分钟数
	FreeMinute             int64 `json:"freeMinute"`             // 免费分钟数
}

type CurrentListenerChatPriceResult struct {
	TextChatPrice        int64  `json:"textChatPrice"`        //文字聊天单价格
	VoiceChatPrice       int64  `json:"voiceChatPrice"`       //语音聊天单价格
	TextChatActualPrice  int64  `json:"textChatActualPrice"`  //优惠后的文字单价格
	VoiceChatActualPrice int64  `json:"voiceChatActualPrice"` //优惠后的语音单价格
	NightAddFlag         string `json:"nightAddFlag"`         // 夜间服务费加收提示
	FreeFlag             string `json:"freeFlag"`             // 优惠标志
	ChatUnitMinute       int64  `json:"chatUnitMinute"`       // 单位时长的分钟数
	FreeMinute           int64  `json:"freeMinute"`           // 免费分钟数
	NewUserDiscount      int64  `json:"newUserDiscount"`      // 新用户减免比率(千分之几)
}

type CurrentShareArg struct {
	TaxAmount            int64 `json:"taxAmount"`            // 税费
	ActualAmount         int64 `json:"actualAmount"`         // 实际总费用
	ShareRateStep1Star5  int64 `json:"shareRateStep1Star5"`  // 1阶段满意评价平台抽佣比率（千分）
	ShareRateStep1Star3  int64 `json:"shareRateStep1Star3"`  // 1阶段一般评价平台抽佣比率（千分）
	ShareRateStep1Star1  int64 `json:"shareRateStep1Star1"`  // 1阶段不满意评价平台抽佣比率（千分）
	ShareAmountStep1Unit int64 `json:"shareAmountStep1Unit"` // 1阶段时长单位个数
	ShareRateStep2Star5  int64 `json:"shareRateStep2Star5"`  // 2阶段满意评价平台抽佣比率（千分）
	ShareRateStep2Star3  int64 `json:"shareRateStep2Star3"`  // 2阶段一般评价平台抽佣比率（千分）
	ShareRateStep2Star1  int64 `json:"shareRateStep2Star1"`  // 2阶段不满意评价平台抽佣比率（千分）
}

func GetCurrentListenerChatPrice(req *CurrentListenerChatPriceArg) *CurrentListenerChatPriceResult {
	var resp CurrentListenerChatPriceResult
	resp.ChatUnitMinute = req.ChatUnitMinute
	resp.NewUserDiscount = req.NewUserDiscount

	// 夜间
	nightAddAmount := CalculateNightAdd(int64(time.Now().Hour()), req.TextChatPrice, req.NightAddPriceRate, DivNumber, req.NightAddPriceHourStart, req.NightAddPriceHourEnd)
	if nightAddAmount > 0 {
		resp.NightAddFlag = fmt.Sprintf("夜间服务，价格已x1.%d费用", req.NightAddPriceRate/10)
	}
	//fmt.Println(nightAddAmount)
	resp.TextChatPrice = req.TextChatPrice + nightAddAmount
	// 文字 税
	resp.TextChatPrice = resp.TextChatPrice + CalculatePrice(resp.TextChatPrice, req.TaxRate, DivNumber)

	// 语音 夜间
	resp.VoiceChatPrice = req.VoiceChatPrice + CalculateNightAdd(int64(time.Now().Hour()), req.VoiceChatPrice, req.NightAddPriceRate, DivNumber, req.NightAddPriceHourStart, req.NightAddPriceHourEnd)
	// 语音 税
	resp.VoiceChatPrice = resp.VoiceChatPrice + CalculatePrice(resp.VoiceChatPrice, req.TaxRate, DivNumber)

	// 打折
	if req.NewUserDiscount == 0 {
		// 取整元
		resp.TextChatPrice = RoundYuan(resp.TextChatPrice)
		resp.TextChatActualPrice = resp.TextChatPrice

		// 取整元
		resp.VoiceChatPrice = RoundYuan(resp.VoiceChatPrice)
		resp.VoiceChatActualPrice = resp.VoiceChatPrice

	} else if req.NewUserDiscount < DivNumber {
		savePrice := RoundDivideInt64(resp.TextChatPrice*(DivNumber-resp.NewUserDiscount), DivNumber)
		resp.TextChatActualPrice = resp.TextChatPrice - savePrice
		// 取整元
		resp.TextChatPrice = RoundYuan(resp.TextChatPrice)
		resp.TextChatActualPrice = RoundYuan(resp.TextChatActualPrice)

		savePrice = RoundDivideInt64(resp.VoiceChatPrice*(DivNumber-resp.NewUserDiscount), DivNumber)
		resp.VoiceChatActualPrice = resp.VoiceChatPrice - savePrice

		//
		resp.VoiceChatPrice = RoundYuan(resp.VoiceChatPrice)
		resp.VoiceChatActualPrice = RoundYuan(resp.VoiceChatActualPrice)

	} else {
		resp.TextChatPrice = RoundYuan(resp.TextChatPrice)
		resp.TextChatActualPrice = resp.TextChatPrice

		resp.VoiceChatPrice = RoundYuan(resp.VoiceChatPrice)
		resp.VoiceChatActualPrice = resp.VoiceChatPrice
	}

	// 优惠标志
	switch req.NewUserDiscount {
	case 0: // 0折 免费
		resp.FreeMinute = req.FreeMinute
		resp.FreeFlag = fmt.Sprintf("XX免费%d分钟", resp.FreeMinute)
	case 500:
		resp.FreeFlag = "XX"
	case 1000:
		resp.FreeFlag = ""
	default:
		resp.FreeFlag = fmt.Sprintf("%.1f折", float64(req.NewUserDiscount/100))
	}
	return &resp
}

// 计算分成
func CalcShareAmount(star, buyUnit, confirmUnit int64, arg *CurrentShareArg) (platformAmount, amount int64) {
	if arg.ActualAmount <= 0 {
		return 0, 0
	}
	canShareAmountSum := arg.ActualAmount - arg.TaxAmount
	if confirmUnit < arg.ShareAmountStep1Unit { // 小于第一阶梯
		// 第一阶梯
		firstStepCanShareAMount := tool.DivideInt64(canShareAmountSum, buyUnit) * (arg.ShareAmountStep1Unit - confirmUnit)
		rate := getShareRate(star, arg.ShareAmountStep1Unit, arg)
		fmt.Println(star, 1, rate)
		platformAmount = CalculatePrice(firstStepCanShareAMount, rate, DivNumber)
		amount = firstStepCanShareAMount - platformAmount

		// 大于第二阶梯的部分
		if buyUnit > (arg.ShareAmountStep1Unit - confirmUnit) {
			rate = getShareRate(star, buyUnit, arg)
			fmt.Println(star, 2, rate)
			platformAmount2 := CalculatePrice(canShareAmountSum-firstStepCanShareAMount, rate, DivNumber)
			amount2 := canShareAmountSum - firstStepCanShareAMount - platformAmount2
			platformAmount += platformAmount2
			amount += amount2
		}
		return
	}
	// 大于等于第一阶梯
	rate := getShareRate(star, confirmUnit+buyUnit, arg)
	fmt.Println(star, 3, rate)
	platformAmount = CalculatePrice(canShareAmountSum, rate, DivNumber)
	amount = canShareAmountSum - platformAmount

	return
}

func getShareRate(star, buyUnit int64, arg *CurrentShareArg) int64 {
	switch star {
	case listenerkey.Rating1Star:
		if buyUnit <= arg.ShareAmountStep1Unit {
			return arg.ShareRateStep1Star1
		}
		return arg.ShareRateStep2Star1

	case listenerkey.Rating3Star:
		if buyUnit <= arg.ShareAmountStep1Unit {
			return arg.ShareRateStep1Star3
		}
		return arg.ShareRateStep2Star3

	case listenerkey.Rating5Star:
		if buyUnit <= arg.ShareAmountStep1Unit {
			return arg.ShareRateStep1Star5
		}
		return arg.ShareRateStep2Star5

	default:
		return 0
	}
}
