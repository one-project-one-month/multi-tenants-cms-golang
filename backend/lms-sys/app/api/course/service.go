package course

import (
	"github.com/multi-tenants-cms-golang/lms-sys/internal/db"
	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
	"github.com/sirupsen/logrus"
)

type CoursesService struct {
	cpb.UnimplementedCourseServiceServer
	store  *db.Store
	logger *logrus.Logger
}

func NewCoursesService(
	store *db.Store,
	logger *logrus.Logger,
) *CoursesService {
	return &CoursesService{
		store:  store,
		logger: logger,
	}
}
