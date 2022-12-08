package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/orderkey"
)

// Crontab 格式
// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/autoClearListenerTodayStatHistoryData.go.
func (l *MqueueScheduler) autoClearListenerTodayStatHistoryData() {
	task := asynq.NewTask(jobtype.ScheduleAutoClearListenerTodayStatHistoryData, nil)
	// every one hour exec
	cronspec := fmt.Sprintf("0 %d * * *", orderkey.ClearListenerTodayStatHistoryDataHour)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【autoClearListenerTodayStatHistoryData】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【autoClearListenerTodayStatHistoryData】 registered an  entry: %q \n", entryID)
	return
}
