package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/listenerkey"
)

// Crontab 格式
// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/autoUpdateListenerUserStat.go.
func (l *MqueueScheduler) autoUpdateListenerUserStat() {
	task := asynq.NewTask(jobtype.ScheduleUpdateListenerUserStat, nil)
	// every one hour exec
	cronspec := fmt.Sprintf("*/%d * * * *", listenerkey.AutoUpdateListenerUserStatIntervalMinute)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【autoUpdateListenerUserStat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【autoUpdateListenerUserStat】 registered an  entry: %q \n", entryID)
	return
}
