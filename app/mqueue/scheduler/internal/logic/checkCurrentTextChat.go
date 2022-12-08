package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/chatkey"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/checkCurrentTextChat.go.
func (l *MqueueScheduler) checkCurrentTextChat() {
	task := asynq.NewTask(jobtype.ScheduleCheckCurrentTextChat, nil)
	// every one minute exec
	cronspec := fmt.Sprintf("*/%d * * * *", chatkey.CheckTextChatUseOutIntervalMinute)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【checkCurrentTextChat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【checkCurrentTextChat】 registered an  entry: %q \n", entryID)
	return
}
