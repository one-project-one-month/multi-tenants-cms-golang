package module

import (
	"context"
	"fmt"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ms *ModulesService) CreateModule(ctx context.Context, req *mpb.CreateModuleRequest) (*mpb.CreateModuleResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	// INFO: To be deleted
	fmt.Printf("rpc_create_module: org FromIncomingContext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	ms.logger.WithFields(logrus.Fields{
		"method":      "CreateModule",
		"module_name": req.ModuleName,
		"org":         orgName,
	}).Info("Creating module")

	if req.ModuleName == "" {
		return nil, status.Error(codes.InvalidArgument, "module name is required")
	}

	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course id is required")
	}

	courseID := module.ConvertStringToGoogleUUID(req.CourseId)

	tenantCheckArgs := repo.IsCourseOwnedByTenantParams{
		CourseID:  courseID,
		Namespace: orgName,
	}

	isOwnedBy, err := ms.store.IsCourseOwnedByTenant(ctx, tenantCheckArgs)
	if err != nil {
		ms.logger.WithError(err).Error("Failed to check if course is owned by tenant")
		return nil, status.Error(codes.Internal, "failed to check course tenant relation")
	}

	if !isOwnedBy {
		return nil, status.Error(codes.PermissionDenied, "not allowed to create a module of a course owned by another tenant")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := db.CreateModuleParams{
		ModuleName:  req.ModuleName,
		CourseID:    courseID,
		Description: module.ConvertStringToText(req.Description),
	}

	m, err := ms.store.CreateModule(dbCtx, args)
	if err != nil {
		ms.logger.WithError(err).Error("Failed to create module")
		return nil, status.Errorf(codes.Internal, "failed to create module: %v", err)
	}

	return &mpb.CreateModuleResponse{
		Module: module.ConvertModuleToProto(m),
	}, nil
}
