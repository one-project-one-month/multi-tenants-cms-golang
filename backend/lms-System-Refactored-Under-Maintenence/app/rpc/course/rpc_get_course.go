package course

import (
	"context"
	coursepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/course"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CourseService) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &coursepb.GetCourseResponse{}
	return resp, nil
}
