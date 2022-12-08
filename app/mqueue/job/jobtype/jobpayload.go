package jobtype

import (
	"time"
)

// DeferCloseChatOrderPayload defer close chat order
type DeferCloseChatOrderPayload struct {
	OrderId string `json:"orderId"`
}

// DeferSendImMsg
type DeferSendImMsgPayload struct {
	KqMsgBuf []byte `json:"kq_msg_buf"` // 发送自定义IM消息kafka消息
}

type DeferCheckChatStatePayload struct {
	Uid         int64 `json:"uid"`
	ListenerUid int64 `json:"listenerUid"`
	OrderType   int64 `json:"orderType"`
}

// PaySuccessNotifyUserPayload pay success notify user TODO
type PaySuccessNotifyUserPayload struct {
	Sn              string    `db:"sn"`                // 订单号
	UserId          int64     `db:"user_id"`           // 下单用户id
	ChatId          int64     `db:"chat_id"`           // 咨询服务id
	Title           string    `db:"title"`             // 标题
	SubTitle        string    `db:"sub_title"`         // 副标题
	Cover           string    `db:"cover"`             // 封面
	Info            string    `db:"info"`              // 介绍
	PeopleNum       int64     `db:"people_num"`        // 容纳人的数量
	RowType         int64     `db:"row_type"`          // 售卖类型0：按房间出售 1:按人次出售
	NeedFood        int64     `db:"need_food"`         // 0:不需要餐食 1:需要参数
	FoodInfo        string    `db:"food_info"`         // 餐食标准
	FoodPrice       int64     `db:"food_price"`        // 餐食价格(分)
	ChatPrice       int64     `db:"chat_price"`        // 咨询服务价格(分)
	MarketChatPrice int64     `db:"market_chat_price"` // 咨询服务市场价格(分)
	ChatBusinessId  int64     `db:"chat_listener_id"`  // 店铺id
	ChatUserId      int64     `db:"chat_user_id"`      // 店铺房东id
	LiveStartDate   time.Time `db:"live_start_date"`   // 开始入住日期
	LiveEndDate     time.Time `db:"live_end_date"`     // 结束入住日期
	LivePeopleNum   int64     `db:"live_people_num"`   // 实际入住人数
	TradeState      int64     `db:"trade_state"`       // -1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期
	TradeCode       string    `db:"trade_code"`        // 确认码
	Remark          string    `db:"remark"`            // 用户下单备注
	OrderTotalPrice int64     `db:"order_total_price"` // 订单总价格（餐食总价格+咨询服务总价格）(分)
	FoodTotalPrice  int64     `db:"food_total_price"`  // 餐食总价格(分)
	ChatTotalPrice  int64     `db:"chat_total_price"`  // 咨询服务总价格(分)
}
