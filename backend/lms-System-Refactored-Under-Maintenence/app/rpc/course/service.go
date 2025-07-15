package course

import (
	coursepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/course"
)

type CoursesService struct {
	coursepb.UnimplementedCourseServiceServer
}

func NewCourseService() *CoursesService {
	return &CoursesService{}
}
