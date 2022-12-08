package timkey

const (
	TimUserSigExpire = 24 * 60 * 60 * 30
)

// 回调命令字
const (
	TimCommandStateChange      = "State.StateChange"            // 状态变更回调
	TimCommandBeforeSendMsg    = "C2C.CallbackBeforeSendMsg"    // 发单聊消息之前回调
	TimCommandAfterSendMsg     = "C2C.CallbackAfterSendMsg"     // 发单聊消息之后回调
	TimCommandAfterMsgReport   = "C2C.CallbackAfterMsgReport"   // 单聊消息已读上报后回调
	TimCommandAfterMsgWithDraw = "C2C.CallbackAfterMsgWithDraw" // 单聊消息撤回后回调
)
