package course

import (
	"context"

	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoursesService) CreateCourse(ctx context.Context, req *cpb.CreateCourseRequest) (*cpb.CreateCourseResponse, error) {
	if req.GetCourseTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "course title is required")
	}
	dbCtx := context.Background()

	_, err := cs.store.CreateNewCourse(dbCtx, req.CourseTitle)
	if err != nil {
		return nil, err
	}

	return &cpb.CreateCourseResponse{
		CourseCreated: &cpb.Course{
			CourseId:    "fasdokn f;ads",
			CourseTitle: "cin sdffdsafads",
			CreatedAt:   nil,
			UpdatedAt:   nil,
		},
	}, status.Errorf(codes.Unimplemented, "method CreateCourse not implemented")
}
