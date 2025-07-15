package module

import (
	"context"
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ModulesService) DeleteModules(ctx context.Context, req *modulepb.DeleteModulesRequest) (*modulepb.DeleteModulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &modulepb.DeleteModulesResponse{}
	return resp, nil
}
