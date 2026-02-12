package module

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ms *ModulesService) DeleteModules(ctx context.Context, req *mpb.DeleteModulesRequest) (*mpb.DeleteModulesResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to retrieve metadata")
	}

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	headerValues := md.Get("x-module-ids")
	if len(headerValues) == 0 {
		return nil, status.Error(codes.InvalidArgument, "module-ids header required")
	}

	var module_ids []string
	if len(headerValues) == 1 {
		module_ids = strings.Split(headerValues[0], ",")
	} else {
		module_ids = headerValues
	}

	ms.logger.WithFields(logrus.Fields{
		"method":     "DeleteModules",
		"module_ids": strings.Join(module_ids, ", "),
		"org":        orgName,
	}).Info("Deleting modules")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, id := range module_ids {
		if id == "" {
			return nil, status.Error(codes.InvalidArgument, "module id is required")
		}

		moduleID := module.ConvertStringToGoogleUUID(id)

		tenantCheckArgs := repo.GetModuleByIDWithTenantParams{
			Namespace: orgName,
			ModuleID:  moduleID,
		}

		m, err := ms.store.GetModuleByIDWithTenant(dbCtx, tenantCheckArgs)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Errorf(codes.NotFound, "module not found")
			}
			return nil, status.Errorf(codes.Internal, "failed to check module existence: %v", err)
		}

		if req.ForceDelete {
			if err := ms.store.DeleteAssociatedQuizzes(dbCtx, moduleID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete quizzes associated to module")
				return nil, status.Errorf(codes.Internal, "failed to force delete quizzes associated to module: %v", err)
			}

			if err := ms.store.DeleteAssociatedLessons(dbCtx, moduleID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete lessons associated to module")
				return nil, status.Errorf(codes.Internal, "failed to force delete lessons associated to module: %v", err)
			}

			if err := ms.store.DeleteModule(dbCtx, moduleID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete module")
				return nil, status.Errorf(codes.Internal, "failed to force delete module: %v", err)
			}

			continue
		}

		hasAssociations, err := ms.store.ModuleHasAssociations(dbCtx, moduleID)
		if err != nil {
			ms.logger.WithError(err).Error("Failed to check associations for module")
			return nil, status.Errorf(
				codes.Internal,
				"failed to check associations for module: %v", err,
			)
		}

		if module.ConvertPgBoolToBool(hasAssociations) {
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"module %s has association records", m.ModuleName,
			)
		}

		if err := ms.store.DeleteModule(dbCtx, moduleID); err != nil {
			ms.logger.WithError(err).Error("Failed to delete module")
			return nil, status.Errorf(codes.Internal, "failed to delete module: %v", err)
		}
	}

	return &mpb.DeleteModulesResponse{
		Message: "Module(s) deleted",
	}, nil
}
