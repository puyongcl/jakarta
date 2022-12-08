package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/orderkey"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/updateListenerLastDayStat.go.
func (l *MqueueScheduler) updateListenerLastDayStat() {
	task := asynq.NewTask(jobtype.ScheduleUpdateListenerLastDayStat, nil)
	//
	cronspec := fmt.Sprintf("10 %d * * *", orderkey.UpdateLastDayListenerStatHour)
	//cronspec := "*/1 * * * *"
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【updateListenerLastDayStat】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【updateListenerLastDayStat】 registered an  entry: %q \n", entryID)
	return
}
