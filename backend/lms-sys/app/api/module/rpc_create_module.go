package module

import (
	"context"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ms *ModulesService) CreateModule(ctx context.Context, req *mpb.CreateModuleRequest) (*mpb.CreateModuleResponse, error) {
	ms.logger.WithFields(logrus.Fields{
		"method":      "CreateModule",
		"module_name": req.ModuleName,
	}).Info("Creating module")

	if req.ModuleName == "" {
		return nil, status.Error(codes.InvalidArgument, "module name is required")
	}

	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course id is required")
	}

	pgCourseId, err := module.ConvertStringToUUID(req.CourseId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format: %v", err)
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := db.CreateModuleParams{
		ModuleName:  req.ModuleName,
		CourseID:    pgCourseId,
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
