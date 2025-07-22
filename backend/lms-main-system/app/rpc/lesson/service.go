package lesson

import (
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	lpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
)

type LessonService struct {
	lpb.UnimplementedLessonServiceServer
	store  db.Store
	logger *logrus.Logger
}

func NewLessonService(store db.Store, logger *logrus.Logger) *LessonService {
	return &LessonService{
		store:  store,
		logger: logger,
	}
}
