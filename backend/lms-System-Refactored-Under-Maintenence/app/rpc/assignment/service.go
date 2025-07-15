package assignment

import (
	assignmentpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignment"
)

type AssignmentsService struct {
	assignmentpb.UnimplementedAssignmentServiceServer
}

func NewAssignmentService() *AssignmentsService {
	return &AssignmentsService{}
}
