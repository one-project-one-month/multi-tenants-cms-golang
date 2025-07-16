package assignment

import (
	"context"
	assignmentpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AssignmentsService) CreateAssignment(ctx context.Context, req *assignmentpb.CreateAssignmentRequest) (*assignmentpb.CreateAssignmentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &assignmentpb.CreateAssignmentResponse{}
	return resp, nil
}
