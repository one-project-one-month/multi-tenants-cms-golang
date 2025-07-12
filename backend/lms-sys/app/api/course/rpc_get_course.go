package course

import (
	"context"
	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
)

func (cs *CoursesService) GetCourse(ctx context.Context, req *cpb.GetCourseRequest) (*cpb.GetCourseResponse, error) {
	return &cpb.GetCourseResponse{
		CourseTitle: req.GetCourseTitle(),
	}, nil
}
