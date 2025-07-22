package module

import (
	"context"
	"fmt"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (ms *ModulesService) ListModules(ctx context.Context, req *emptypb.Empty) (*mpb.ListModulesResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	// INFO: To be deleted
	fmt.Printf("rpc_list_modules: org FromIncomingContext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	ms.logger.WithFields(logrus.Fields{
		"method": "ListModules",
		"org":    orgName,
	}).Info("Listing modules")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modules, err := ms.store.ListModulesWithTenant(dbCtx, orgName)
	if err != nil {
		ms.logger.WithError(err).Errorf("Failed to list module")
		return nil, status.Errorf(codes.Internal, "failed to list module: %v", err)
	}

	return &mpb.ListModulesResponse{
		Modules: module.ConvertModulesToProto(modules),
	}, nil
}
