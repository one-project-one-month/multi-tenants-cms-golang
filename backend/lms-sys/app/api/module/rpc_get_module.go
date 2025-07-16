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

func (ms *ModulesService) GetModule(ctx context.Context, req *mpb.GetModuleRequest) (*mpb.GetModuleResponse, error) {
	ms.logger.WithFields(logrus.Fields{
		"method":    "GetModuleByID",
		"module_id": req.ModuleId,
	}).Info("Getting module")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgUUID, err := module.ConvertStringToUUID(req.ModuleId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid module id format %v", err)
	}

	m, err := ms.store.GetModuleByID(dbCtx, pgUUID)
	if err != nil {
		ms.logger.WithError(err).Errorf("Failed to get module")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "module not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get module: %v", err)
	}

	return &mpb.GetModuleResponse{
		Module: module.ConvertModuleToProto(m),
	}, nil
}
