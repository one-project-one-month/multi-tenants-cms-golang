package module

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ms *ModulesService) DeleteModules(ctx context.Context, req *mpb.DeleteModulesRequest) (*mpb.DeleteModulesResponse, error) {
	ms.logger.WithFields(logrus.Fields{
		"method":     "DeleteModules",
		"module_ids": req.Ids,
	})

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, id := range req.Ids {
		if id == "" {
			return nil, status.Error(codes.InvalidArgument, "module id is required")
		}

		pgUUID, err := module.ConvertStringToUUID(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid module id format: %v", err)
		}

		m, err := ms.store.GetModuleByID(dbCtx, pgUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Errorf(codes.NotFound, "module not found")
			}
			return nil, status.Errorf(codes.Internal, "failed to check module existence: %v", err)
		}

		if req.ForceDelete {
			if err := ms.store.DeleteAssociatedQuizzes(dbCtx, pgUUID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete quizzes associated to module")
				return nil, status.Errorf(codes.Internal, "failed to force delete quizzes associated to module: %v", err)
			}

			if err := ms.store.DeleteAssociatedLessons(dbCtx, pgUUID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete lessons associated to module")
				return nil, status.Errorf(codes.Internal, "failed to force delete lessons associated to module: %v", err)
			}

			if err := ms.store.DeleteModule(dbCtx, pgUUID); err != nil {
				ms.logger.WithError(err).Error("Failed to force delete module")
				return nil, status.Errorf(codes.Internal, "failed to force delete module: %v", err)
			}

			continue
		}

		hasAssociations, err := ms.store.ModuleHasAssociations(dbCtx, pgUUID)
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

		if err := ms.store.DeleteModule(dbCtx, pgUUID); err != nil {
			ms.logger.WithError(err).Error("Failed to delete module")
			return nil, status.Errorf(codes.Internal, "failed to delete module: %v", err)
		}
	}

	return &mpb.DeleteModulesResponse{
		Message: "Module(s) deleted",
	}, nil
}
