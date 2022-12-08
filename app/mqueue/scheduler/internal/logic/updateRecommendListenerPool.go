package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
)

// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/updateNewUserRecommendListenerEveryDay.go.
func (l *MqueueScheduler) updateRecommendListenerPool() {
	task := asynq.NewTask(jobtype.ScheduleUpdateRecommendListenerPool, nil)
	//
	//cronspec := fmt.Sprintf("10 %d * * *", orderkey.UpdateLastDayListenerStatHour)
	cronspec := "10 0 * * *"
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【updateRecommendListenerPool】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【updateRecommendListenerPool】 registered an  entry: %q \n", entryID)
	return
}
