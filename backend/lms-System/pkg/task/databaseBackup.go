package task

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type DatabaseBackup struct {
	DbUrl      string
	BackUpPath string
}

func NewDatabaseBackupTask(dbUrl string, backupPath string) (*asynq.Task, error) {
	payload, err := json.Marshal(&DatabaseBackup{
		DbUrl:      dbUrl,
		BackUpPath: backupPath,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDatabaseBackup, payload), nil
}
func ScheduleMonthlyBackup(dbUrl, backupPath string) (*asynq.Task, *asynq.Scheduler, error) {
	payload, err := json.Marshal(DatabaseBackup{
		DbUrl:      dbUrl,
		BackUpPath: backupPath,
	})
	if err != nil {
		return nil, nil, err
	}

	task := asynq.NewTask(TypeDatabaseBackup, payload)

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		&asynq.SchedulerOpts{},
	)

	return task, scheduler, nil
}
