package module

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ms *ModulesService) UpdateModule(ctx context.Context, req *mpb.UpdateModuleRequest) (*mpb.UpdateModuleResponse, error) {
	ms.logger.WithFields(logrus.Fields{
		"method":    "UpdateModule",
		"module_id": req.ModuleId,
	}).Info("Updating module")

	if req.ModuleId == "" {
		return nil, status.Error(codes.InvalidArgument, "module id is required")
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modulepgUUID, err := module.ConvertStringToUUID(req.ModuleId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid module id format: %v", err)
	}

	_, err = ms.store.GetModuleByID(dbCtx, modulepgUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "module not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to check module existense: %v", err)
	}

	coursepgUUID, err := module.ConvertStringToUUID(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid course id format: %v", err)
	}

	args := db.UpdateModuleByIDParams{
		ModuleID:    modulepgUUID,
		ModuleName:  req.ModuleName,
		CourseID:    coursepgUUID,
		Description: module.ConvertStringToText(req.Description),
	}

	m, err := ms.store.UpdateModuleByID(dbCtx, args)
	if err != nil {
		ms.logger.WithError(err).Error("Failed to update module")
		return nil, status.Errorf(codes.Internal, "failed to update module: %v", err)
	}

	return &mpb.UpdateModuleResponse{
		Module: module.ConvertModuleToProto(m),
	}, nil
}
