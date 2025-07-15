package courses

import (
	"context"
	coursespb "github.com/multi-tenants-cms-golang/lms-system/protogen/courses"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoursesService) CreateCourse(ctx context.Context, req *coursespb.CreateCourseRequest) (*coursespb.CreateCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &coursespb.CreateCourseResponse{}
	return resp, nil
}
