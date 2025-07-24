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

func (es *EnrollmentService) CreateEnrollment(
	ctx context.Context,
	req *epb.CreateEnrollmentRequest,
) (*epb.CreateEnrollmentResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to retrieve metadata")
	}

	// TODO: delete this
	fmt.Printf("rpc_create_enrollment: metadata fromincomingcontext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	es.logger.WithFields(logrus.Fields{
		"method":     "CreateEnrollment",
		"course_id":  req.CourseId,
		"student_id": req.StudentId,
		"org":        orgName,
	}).Info("Creating new enrollment")

	if req.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student id is required")
	}

	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course id is required")
	}

	courseID := enrollment.ConvertStringToGoogleUUID(req.CourseId)

	tenantCheckArgs := db.IsCourseOwnedByTenantParams{
		CourseID:  courseID,
		Namespace: orgName,
	}
	isOwnedByTenant, err := es.store.IsCourseOwnedByTenant(ctx, tenantCheckArgs)
	if err != nil {
		es.logger.WithError(err).Error("Failed to check if course is owned by tenant")
		return nil, status.Error(codes.Internal, "failed to check course tenant relation")
	}

	if !isOwnedByTenant {
		return nil, status.Error(codes.PermissionDenied, "not allowed to enroll to a course owned by another tenant")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	studentID := enrollment.ConvertStringToGoogleUUID(req.StudentId)

	studentCheckArgs := db.IsEnrolleeStudentParams{
		Namespace: orgName,
		LmsUserID: studentID,
	}

	isEnrolleeStudent, err := es.store.IsEnrolleeStudent(dbCtx, studentCheckArgs)
	if err != nil {
		es.logger.WithError(err).Errorf("Failed to check if enrollee is student\n student id: %s\ncourse id: %s\n",
			req.StudentId,
			req.CourseId,
		)
		return nil, status.Error(codes.Internal, "failed to check if enrollee is student")
	}

	if !isEnrolleeStudent {
		return nil, status.Error(codes.PermissionDenied, "enrollee is not a student")
	}

	isEnrollmentExist, err := es.store.IsEnrollmentExist(dbCtx, db.IsEnrollmentExistParams{
		StudentID: studentID,
		CourseID:  courseID,
	})
	if err != nil {
		es.logger.WithError(err).Errorf("Failed to check if enrollment exists\n student id: %s\ncourse id: %s\n",
			req.StudentId,
			req.CourseId,
		)
		return nil, status.Error(codes.Internal, "failed to check if enrollment already exists")
	}

	if isEnrollmentExist {
		return nil, status.Error(codes.AlreadyExists, "student have already enrolled the course")
	}

	defaultDueDate := time.Now().Add(30 * 24 * time.Hour)

	createArgs := db.CreateEnrollmentParams{
		StudentID: studentID,
		CourseID:  courseID,
		DueDate:   enrollment.ConvertTimeToPgTime(defaultDueDate),
		Status:    "Pending",
	}

	enrollmentID, err := es.store.CreateEnrollment(dbCtx, createArgs)
	if err != nil {
		es.logger.WithError(err).Errorf("Failed to create a new enrollment\n student id: %s\ncourse id: %s\n",
			req.StudentId,
			req.CourseId,
		)
		return nil, status.Error(codes.Internal, "failed to create a new enrollment")
	}

	detailsArgs := db.GetEnrollmentByIDAndTenantParams{
		Namespace:    orgName,
		EnrollmentID: enrollmentID,
	}

	enrollmentDetails, err := es.store.GetEnrollmentByIDAndTenant(dbCtx, detailsArgs)
	if err != nil {
		es.logger.WithError(err).Errorf("Failed to retrieve the created enrollment\nenrollment id: %s\n", enrollmentID)
		return nil, status.Error(codes.Internal, "failed to retrieve created enrollment")
	}

	return &epb.CreateEnrollmentResponse{
		Enrollment: enrollment.ConvertEnrollmentDetailsToProto(enrollmentDetails),
	}, nil
}
