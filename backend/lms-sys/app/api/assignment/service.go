package assignment

import (
	"github.com/multi-tenants-cms-golang/lms-sys/internal/db"
	atp "github.com/multi-tenants-cms-golang/lms-sys/protogen/assignments"
	"github.com/sirupsen/logrus"
)

type AssignmentsService struct {
	atp.UnimplementedAssignmentServiceServer

	store  *db.Store
	logger *logrus.Logger
}

func NewAssignmentsService(
	store *db.Store,
	logger *logrus.Logger,
) *AssignmentsService {
	return &AssignmentsService{
		store: store,
		logger: logger,
	}
}
