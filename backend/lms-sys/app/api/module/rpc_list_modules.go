package module

import (
	"context"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/module"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (ms *ModuleService) ListModules(ctx context.Context, req *emptypb.Empty) (*mpb.ListModulesResponse, error) {
	ms.logger.WithFields(logrus.Fields{
		"method": "ListModules",
	}).Info("Listing module")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modules, err := ms.store.ListModules(dbCtx)
	if err != nil {
		ms.logger.WithError(err).Errorf("Failed to list module")
		return nil, status.Errorf(codes.Internal, "failed to list module: %v", err)
	}

	return &mpb.ListModulesResponse{
		Modules: module.ConvertModulesToProto(modules),
	}, nil
}
