package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/chatkey"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/checkCurrentVoiceChat.go.
func (l *MqueueScheduler) checkCurrentVoiceChat() {
	task := asynq.NewTask(jobtype.ScheduleCheckCurrentVoiceChat, nil)
	// every one minute exec
	cronspec := fmt.Sprintf("@every %ds", chatkey.CheckVoiceChatUseOutIntervalSecond)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【checkCurrentVoiceChat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【checkCurrentVoiceChat】 registered an  entry: %q \n", entryID)
	return
}
