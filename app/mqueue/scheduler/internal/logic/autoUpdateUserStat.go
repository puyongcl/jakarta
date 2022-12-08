package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
)

// Crontab 格式
// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/autoUpdateUserStat.go.
func (l *MqueueScheduler) autoUpdateUserStat() {
	task := asynq.NewTask(jobtype.ScheduleAutoUpdateUserStat, nil)
	// every one hour exec
	cronspec := fmt.Sprintf("*/%d * * * *", 1)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【autoUpdateUserStat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【autoUpdateUserStat】 registered an  entry: %q \n", entryID)
	return
}
