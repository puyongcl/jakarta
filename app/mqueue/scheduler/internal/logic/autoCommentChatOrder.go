package logic

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/common/key/orderkey"
)

// Crontab 格式
// scheduler job ------> go-zero-jakarta/app/mqueue/job/internal/logic/autoCommentChatOrder.go.
func (l *MqueueScheduler) autoCommentChatOrder() {
	task := asynq.NewTask(jobtype.ScheduleAutoCommentAndFinishChatOrder, nil)
	// every one hour exec
	cronspec := fmt.Sprintf("0 */%d * * *", orderkey.AutoGoodCommentScheduleIntervalHour)
	entryID, err := l.svcCtx.Scheduler.Register(cronspec, task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【autoCommentChatOrder】 registered  err:%+v , task:%+v", err, task)
		return
	}
	fmt.Printf("【autoCommentChatOrder】 registered an  entry: %q \n", entryID)
	return
}
