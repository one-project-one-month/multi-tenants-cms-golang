package module

import (
	"context"
	"database/sql"
	"errors"
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

func (ms *ModulesService) UpdateModule(ctx context.Context, req *mpb.UpdateModuleRequest) (*mpb.UpdateModuleResponse, error) {
	// TODO: Verify whether to use tenant filtering for updates
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	// INFO: To be deleted
	fmt.Printf("rpc_update_modules: org FromIncomingContext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	ms.logger.WithFields(logrus.Fields{
		"method":    "UpdateModule",
		"module_id": req.ModuleId,
		"org":       orgName,
	}).Info("Updating module")

	if req.ModuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "module id is required")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moduleID := module.ConvertStringToGoogleUUID(req.ModuleId)

	args := repo.GetModuleByIDWithTenantParams{
		ModuleID:  moduleID,
		Namespace: orgName,
	}

	currentModule, err := ms.store.GetModuleByIDWithTenant(dbCtx, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "module not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to check module existence: %v", err)
	}

	moduleName := currentModule.ModuleName
	if req.ModuleName != nil {
		moduleName = *req.ModuleName
	}

	courseID := currentModule.CourseID
	if req.CourseId != nil && *req.CourseId != currentModule.CourseID.String() {
		// INFO: Check if new parent course is under the same tenant
		args := repo.IsCourseOwnedByTenantParams{
			CourseID:  module.ConvertStringToGoogleUUID(*req.CourseId),
			Namespace: orgName,
		}

		isOwnedByCurrentTenant, err := ms.store.IsCourseOwnedByTenant(ctx, args)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Error(codes.NotFound, "new parent course not found")
			}
			return nil, status.Errorf(codes.Internal, "failed to check if the new course is owned by the same tenant: %v", err)
		}

		if !isOwnedByCurrentTenant {
			return nil, status.Errorf(codes.PermissionDenied, "new course is not under the same tenant")
		}

		courseID = module.ConvertStringToGoogleUUID(*req.CourseId)
	}

	moduleDescription := currentModule.Description
	if req.Description != nil {
		moduleDescription = module.ConvertStringToText(*req.Description)
	}

	// INFO: Normal update
	updateArgs := db.UpdateModuleByIDParams{
		ModuleID:    moduleID,
		ModuleName:  moduleName,
		CourseID:    courseID,
		Description: moduleDescription,
	}

	m, err := ms.store.UpdateModuleByID(dbCtx, updateArgs)
	if err != nil {
		ms.logger.WithError(err).Error("Failed to update module")
		return nil, status.Errorf(codes.Internal, "failed to update module: %v", err)
	}

	return &mpb.UpdateModuleResponse{
		Module: module.ConvertModuleToProto(m),
	}, nil
}
