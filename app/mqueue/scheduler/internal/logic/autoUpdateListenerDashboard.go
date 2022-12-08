package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/listenerkey"
)

// Crontab 格式
// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/autoUpdateListenerDashboard.go.
func (l *MqueueScheduler) autoUpdateListenerDashboard() {
	task := asynq.NewTask(jobtype.ScheduleUpdateListenerDashboardStat, nil)
	// every one hour exec
	cronspec := fmt.Sprintf("*/%d * * * *", listenerkey.AutoUpdateListenerUserStatIntervalMinute+1)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【autoUpdateListenerDashboard】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【autoUpdateListenerDashboard】 registered an  entry: %q \n", entryID)
	return
}
