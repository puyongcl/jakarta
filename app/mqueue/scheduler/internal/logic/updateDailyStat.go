package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/updateListenerLastDayStat.go.
func (l *MqueueScheduler) updateDailyStat() {
	task := asynq.NewTask(jobtype.ScheduleUpdateDailyStat, nil)
	//
	cronspec := "1 0 * * *"
	//cronspec := "*/1 * * * *"
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【updateDailyStat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【updateDailyStat】 registered an  entry: %q \n", entryID)
	return
}
