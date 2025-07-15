package course

import (
	coursepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/course"
)

type CourseService struct {
	coursepb.UnimplementedCourseServiceServer
}

func NewCourseService() *CourseService {
	return &CourseService{}
}
