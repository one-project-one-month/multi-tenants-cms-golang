package module

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (ms *ModulesService) GetModule(ctx context.Context, req *mpb.GetModuleRequest) (*mpb.GetModuleResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	// INFO: To be deleted
	fmt.Printf("rpc_get_module: org FromIncomingContext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	ms.logger.WithFields(logrus.Fields{
		"method":    "GetModuleByID",
		"module_id": req.ModuleId,
		"org":       orgName,
	}).Info("Getting module")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := repo.GetModuleByIDWithTenantParams{
		ModuleID:  module.ConvertStringToGoogleUUID(req.ModuleId),
		Namespace: orgName,
	}

	m, err := ms.store.GetModuleByIDWithTenant(dbCtx, args)
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
