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
	"google.golang.org/protobuf/types/known/emptypb"
)

func (es *EnrollmentService) ListEnrollments(
	ctx context.Context,
	req *emptypb.Empty,
) (*epb.ListEnrollmentsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to retrieve metadata")
	}

	// TODO: delete this
	fmt.Printf("rpc_list_enrollments: metadata fromincomingcontext: %+v\n", md)

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

	es.logger.WithFields(logrus.Fields{
		"method":       "ListEnrollments",
		"requestor_id": requestor_id,
		"org":          orgName,
	}).Info("Listing enrollments")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var enrollments []db.EnrollmentDetails

	// INFO: requestor is admin
	roleCheckArgs := db.IsUserInRoleParams{
		Namespace:   orgName,
		LmsUserID:   requestor_id,
		LmsRoleName: "ADMIN",
	}
	isAdmin, err := es.store.IsUserInRole(dbCtx, roleCheckArgs)
	if err != nil {
		es.logger.WithError(err).Error("Failed to check requestor's role")
		return nil, status.Error(codes.Internal, "failed to check requestor's role")
	}
	if isAdmin {
		enrollments, err = es.store.ListEnrollmentsByAdminFilters(dbCtx, db.ListEnrollmentsByAdminFiltersParams{
			Namespace: orgName,
		})
		if err != nil {
			es.logger.WithError(err).Error("Failed to list enrollments for admin")
			return nil, status.Error(codes.Internal, "failed to list enrollments for admin")
		}
		return &epb.ListEnrollmentsResponse{
			Enrollments: enrollment.ConvertEnrollmentDetailsListToProto(enrollments),
		}, nil
	}

	// INFO: requestor is instructor
	roleCheckArgs.LmsRoleName = "INSTRUCTOR"
	isInstructor, err := es.store.IsUserInRole(dbCtx, roleCheckArgs)
	if err != nil {
		es.logger.WithError(err).Error("Failed to check requestor's role")
		return nil, status.Error(codes.Internal, "failed to check requestor's role")
	}
	if isInstructor {
		enrollments, err = es.store.ListEnrollmentsByInstructorFilters(dbCtx, db.ListEnrollmentsByInstructorFiltersParams{
			Namespace:    orgName,
			InstructorID: requestor_id,
		})
		if err != nil {
			es.logger.WithError(err).Error("Failed to list enrollments for admin")
			return nil, status.Error(codes.Internal, "failed to list enrollments for admin")
		}
		return &epb.ListEnrollmentsResponse{
			Enrollments: enrollment.ConvertEnrollmentDetailsListToProto(enrollments),
		}, nil
	}

	return nil, status.Error(codes.PermissionDenied, "you are not allowed to view enrollments")
}
