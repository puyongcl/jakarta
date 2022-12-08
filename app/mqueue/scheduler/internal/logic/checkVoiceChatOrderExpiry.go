package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/orderkey"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/checkVoiceChatOrderExpiry.go.
func (l *MqueueScheduler) checkVoiceChatOrderExpiry() {
	task := asynq.NewTask(jobtype.ScheduleCheckVoiceChatOrderExpiry, nil)
	// every one minute exec
	cronspec := fmt.Sprintf("*/%d * * * *", orderkey.CheckOrderExpiryIntervalMinutes)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【checkVoiceChatOrderExpiry】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【checkVoiceChatOrderExpiry】 registered an  entry: %q \n", entryID)
	return
}
