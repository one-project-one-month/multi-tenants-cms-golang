package enrollment

import (
	"context"
	"fmt"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/enrollment"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	epb "github.com/multi-tenants-cms-golang/lms-sys/protogen/enrollments"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (es *EnrollmentService) DeleteEnrollment(
	ctx context.Context,
	req *epb.DeleteEnrollmentRequest,
) (*epb.DeleteEnrollmentResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to retrieve metadata")
	}

	// TODO: delete this
	fmt.Printf("rpc_delete_enrollment: metadata fromincomingcontext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	IDs := md.Get("x-requestor-id")

	if len(IDs) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "missing requestor id")
	}

	requestor_id := enrollment.ConvertStringToGoogleUUID(IDs[0])

	enrollment_id := enrollment.ConvertStringToGoogleUUID(req.EnrollmentId)

	es.logger.WithFields(logrus.Fields{
		"method":        "ListEnrollments",
		"requestor_id":  requestor_id,
		"enrollment_id": enrollment_id,
		"org":           orgName,
	}).Info("Deleting enrollment")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isEnrollmentExists, err := es.store.IsEnrollmentExists(dbCtx, enrollment_id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check enrollment's existence")
	}

	if !isEnrollmentExists {
		return nil, status.Error(codes.NotFound, "enrollment to delete doesn't exist")
	}

	studentCheckArgs := db.IsUserInRoleParams{
		Namespace:   orgName,
		LmsUserID:   requestor_id,
		LmsRoleName: "STUDENT",
	}

	isStudent, err := es.store.IsUserInRole(dbCtx, studentCheckArgs)
	if err != nil {
		es.logger.WithError(err).Error("Failed to check requestor is student")
		return nil, status.Error(codes.Internal, "failed to check requestor is student")
	}

	if !isStudent {
		return nil, status.Error(codes.PermissionDenied, "you are not allowed to delete the enrollment")
	}

	isBelong, err := es.store.IsEnrollmentBelongsToStudent(dbCtx, db.IsEnrollmentBelongsToStudentParams{
		EnrollmentID: enrollment_id,
		StudentID:    requestor_id,
	})
	if err != nil {
		es.logger.WithError(err).Errorf("Failed to check if the enrollment '%s' belongs to student '%s'", enrollment_id, requestor_id)
		return nil, status.Error(codes.Internal, "failed to check if the enrollment belongs to student")
	}

	if !isBelong {
		return nil, status.Error(codes.NotFound, "not found the requested enrollment")
	}

	if err := es.store.DeleteEnrollmentByID(dbCtx, enrollment_id); err != nil {
		es.logger.WithError(err).Errorf("Failed to delete enrollment id: '%s'", enrollment_id)
		return nil, status.Error(codes.Internal, "failed to delete enrollment")
	}

	return &epb.DeleteEnrollmentResponse{
		Message: "Successfully deleted the enrollment",
	}, nil
}
