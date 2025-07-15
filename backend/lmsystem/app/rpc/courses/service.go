package courses

import (
	coursespb "github.com/multi-tenants-cms-golang/lms-system/protogen/courses"
)

type CoursesService struct {
	coursespb.UnimplementedCoursesServiceServer
}

func NewCoursesService() *CoursesService {
	return &CoursesService{}
}
