package enrollment

import (
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	epb "github.com/multi-tenants-cms-golang/lms-sys/protogen/enrollments"
	"github.com/sirupsen/logrus"
)

type EnrollmentService struct {
	epb.UnimplementedEnrollmentServiceServer
	store  db.Store
	logger *logrus.Logger
}

func NewEnrollmentService(store db.Store, logger *logrus.Logger) *EnrollmentService {
	return &EnrollmentService{
		store:  store,
		logger: logger,
	}
}
