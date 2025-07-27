package quiz

import (
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	qpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/quiz"
	"github.com/sirupsen/logrus"
)

type QuizService struct {
	store  db.Store
	logger *logrus.Logger
	qpb.UnimplementedQuizServiceServer
}

func NewQuizService(store db.Store, logger *logrus.Logger) *QuizService {
	return &QuizService{
		store:  store,
		logger: logger,
	}
}
