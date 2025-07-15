package module

import (
	"context"
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ModulesService) UpdateModule(ctx context.Context, req *modulepb.UpdateModuleRequest) (*modulepb.UpdateModuleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &modulepb.UpdateModuleResponse{}
	return resp, nil
}
