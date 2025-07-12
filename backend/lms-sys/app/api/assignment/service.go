package assignment

import (
	database "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	atp "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignments"
	"github.com/sirupsen/logrus"
)

type AssignmentsService struct {
	atp.UnimplementedAssignmentServiceServer

	store  *database.Store
	logger *logrus.Logger
}

func NewAssignmentsService(
	store *database.Store,
	logger *logrus.Logger,
) *AssignmentsService {
	return &AssignmentsService{
		store:  store,
		logger: logger,
	}
}
