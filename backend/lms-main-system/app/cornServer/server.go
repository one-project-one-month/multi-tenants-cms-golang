package cornServer

import (
	"github.com/hibiken/asynq"
	"github.com/multi-tenants-cms-golang/lms-sys/app/cornServer/handler"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/task"
	"github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	logger    *logrus.Logger
	redisAddr string
	handlers  *handler.CornHandler
	dbUrl     string
}

func NewServer(
	logger *logrus.Logger,
	redisAddr string,
	handler *handler.CornHandler,
	dbUrl string,
) *Server {
	return &Server{
		logger:    logger,
		redisAddr: redisAddr,
		handlers:  handler,
		dbUrl:     dbUrl,
	}
}
func (s *Server) RunServer() error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: s.redisAddr,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: s.redisAddr},
		&asynq.SchedulerOpts{
			Logger:   s.logger,
			Location: time.UTC,
		},
	)
	backupTask, err := task.NewDatabaseBackupTask(s.dbUrl, "/Users/swanhtet/Desktop/multi-tenants-cms-golang/backend/lms-main-system/backups")
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Failed to create backups task")
	}
	_, err = scheduler.Register("@every 1m", backupTask)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err,
		}).Warnf("Failed to register scheduler")
	}
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeDatabaseBackup, s.handlers.HandleDatabaseBackUp)

	go func() {
		if err := scheduler.Run(); err != nil {
			s.logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Failed to run scheduler")
		}
	}()
	s.logger.Info("Task server is listening")
	if err := srv.Run(mux); err != nil {
		s.logger.Info("could not run the task server")
		s.logger.Fatal(err.Error())
		return err
	}
	return nil
}
