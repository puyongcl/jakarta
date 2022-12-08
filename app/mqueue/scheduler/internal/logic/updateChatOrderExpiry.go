package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/orderkey"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/updateChatOrderExpiry.go.
func (l *MqueueScheduler) updateChatOrderExpiry() {
	task := asynq.NewTask(jobtype.ScheduleUpdateVoiceChatOrderExpiry, nil)
	// every one minute exec
	cronspec := fmt.Sprintf("*/%d * * * *", orderkey.UpdateOrderExpiryIntervalMinutes)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【updateChatOrderExpiry】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【updateChatOrderExpiry】 registered an  entry: %q \n", entryID)
	return
}
