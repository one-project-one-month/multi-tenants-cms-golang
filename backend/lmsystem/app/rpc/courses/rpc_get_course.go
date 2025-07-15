package courses

import (
	"context"
	coursespb "github.com/multi-tenants-cms-golang/lms-system/protogen/courses"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoursesService) GetCourse(ctx context.Context, req *coursespb.GetCourseRequest) (*coursespb.GetCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &coursespb.GetCourseResponse{}
	return resp, nil
}
