package assignment

import (
	assignmentpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignment"
)

type AssignmentService struct {
	assignmentpb.UnimplementedAssignmentServiceServer
}

func NewAssignmentService() *AssignmentService {
	return &AssignmentService{}
}
